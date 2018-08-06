package main

import (
	"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	client "github.com/hyperledger/sawtooth-sdk-go/protobuf/client_event_pb2"
	events "github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	validator "github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	zmq "github.com/pebbe/zmq4"
	addresser "github.com/rico-bee/leopark/address"
	"log"
)

const (
	NULL_BLOCK_ID = "0000000000000000"
)

type Subscriber struct {
	connection messaging.Connection
	db         *DbServer
	events     chan *events.EventList
	done       <-chan interface{}
}

func (s *Subscriber) checkResponse(corId string, message proto.Message) (validator.Message_MessageType, error) {
	_, msg, err := s.connection.RecvMsgWithId(corId)
	log.Println("received message:" + corId)
	if err != nil {
		log.Println(err.Error())
		return validator.Message_DEFAULT, err
	}

	err = proto.Unmarshal(msg.Content, message)
	if err != nil {
		return validator.Message_DEFAULT, err
	}
	return msg.MessageType, nil
}

func (s *Subscriber) Start(knownIds []string) error {
	if len(knownIds) == 0 {
		knownIds = []string{NULL_BLOCK_ID}
	}

	for _, id := range knownIds {
		log.Println("known id:" + id)
	}
	log.Println("subscribing to state delta events")
	blockSub := &events.EventSubscription{
		EventType: "sawtooth/block-commit",
	}
	deltaSub := &events.EventSubscription{
		EventType: "sawtooth/state-delta",
		Filters: []*events.EventFilter{
			&events.EventFilter{
				Key:         "address",
				MatchString: "^" + addresser.NS + ".*",
				FilterType:  events.EventFilter_REGEX_ANY,
			},
		},
	}
	req := &client.ClientEventsSubscribeRequest{
		LastKnownBlockIds: knownIds,
		Subscriptions:     []*events.EventSubscription{blockSub, deltaSub},
	}
	rawBytes, err := proto.Marshal(req)
	if err != nil {
		log.Println("corrupted msg" + err.Error())
	}
	corId, err := s.connection.SendNewMsg(validator.Message_CLIENT_EVENTS_SUBSCRIBE_REQUEST, rawBytes)
	log.Println("send subscribe request:" + corId)
	if err != nil {
		log.Fatal("failed to send subscribe request:" + err.Error())
	}
	msg := &client.ClientEventsSubscribeResponse{}
	msgType, err := s.checkResponse(corId, msg)
	if msgType == validator.Message_CLIENT_EVENTS_SUBSCRIBE_RESPONSE {
		if msg.Status == client.ClientEventsSubscribeResponse_UNKNOWN_BLOCK {
			// retrying the service...
			log.Println("retrying to subscribe...")
			s.Start([]string{})
		}
	} else {
		return errors.New("failed to subscribe to validator")
	}

	for {
		es := s.subscribe()
		processEventList(s.done, es, s.db)
	}
}

func (s *Subscriber) subscribe() <-chan *events.Event {
	log.Println("tried to subscribe")
	_, msg, err := s.connection.RecvMsg()
	eventStream := make(chan *events.Event)
	go func() {
		defer close(eventStream)
		if err != nil {
			log.Println("failed to receive the message:" + err.Error())
		} else {
			log.Println("message received:" + msg.MessageType.String())
			var eventsList events.EventList
			err = proto.Unmarshal(msg.Content, &eventsList)
			if err != nil {
				log.Println("err:" + err.Error())
				return
			}
			for _, e := range eventsList.Events {
				select {
				case <-s.done:
					return
				case eventStream <- e:
				}
			}
		}
	}()
	return eventStream
}

func NewSubscriber(validatorUrl string, done <-chan interface{}, db *DbServer) *Subscriber {
	zmqCtx, err := zmq.NewContext()
	if err != nil {
		log.Fatalln("cannot create zmq context")
	}

	conn, err := messaging.NewConnection(zmqCtx, zmq.DEALER, validatorUrl, false)
	if err != nil {
		log.Println("failed to create connection")
	}

	return &Subscriber{
		connection: conn,
		events:     make(chan *events.EventList),
		done:       done,
		db:         db,
	}
}

func (s *Subscriber) Stop() {
	req := &client.ClientEventsUnsubscribeRequest{}
	rawBytes, err := proto.Marshal(req)
	if err != nil {
		log.Println("corrupted msg" + err.Error())
	}
	corId, err := s.connection.SendNewMsg(validator.Message_CLIENT_EVENTS_UNSUBSCRIBE_REQUEST, rawBytes)
	msg := &client.ClientEventsUnsubscribeResponse{}
	msgType, err := s.checkResponse(corId, msg)
	if msgType == validator.Message_CLIENT_EVENTS_SUBSCRIBE_RESPONSE {
		if msg.Status != client.ClientEventsUnsubscribeResponse_OK {
			log.Println("failed to unsubscribe from validator" + msg.Status.String())
		}
	} else {
		log.Fatal("corrupted response from:Message_CLIENT_EVENTS_UNSUBSCRIBE_REQUEST")
	}
}
