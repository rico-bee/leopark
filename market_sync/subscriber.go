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
	done       <-chan interface{}
}

func (s *Subscriber) checkResponse(corId string) (validator.Message_MessageType, proto.Message, error) {
	id, msg, err := s.connection.RecvMsgWithId(corId)

	log.Println("received message:" + id)
	if err != nil {
		log.Println(err.Error())
		return validator.Message_DEFAULT, nil, err
	}

	var message proto.Message
	err = proto.Unmarshal(msg.Content, message)
	if err != nil {
		return validator.Message_DEFAULT, nil, err
	}
	return validator.Message_DEFAULT, nil, err
}

func (s *Subscriber) Start(knownIds []string) error {
	if len(knownIds) == 0 {
		knownIds = []string{NULL_BLOCK_ID}
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
	msgType, msg, err := s.checkResponse(corId)
	if msgType == validator.Message_CLIENT_EVENTS_SUBSCRIBE_RESPONSE {
		subscribeRes, ok := msg.(*client.ClientEventsSubscribeResponse)
		if !ok {
			log.Println("failed to parse out the reponse")
			return errors.New("cannot parse out subscribe response")
		}
		if subscribeRes.Status == client.ClientEventsSubscribeResponse_UNKNOWN_BLOCK {
			// retrying the service...
			s.Start([]string{})
		}
	} else {
		return errors.New("failed to subscribe to validator")
	}

	eventStream := make(chan *events.EventList)
	defer close(eventStream)
	go processEventList(eventStream, s.db)
	for {
		select {
		case <-s.done:
			log.Println("gracefully exit subscribing...")
			return nil
		default:
		}
		s.subscribe(eventStream)
	}
}

func (s *Subscriber) subscribe(eventStream chan<- *events.EventList) {
	id, msg, err := s.connection.RecvMsg()
	if err != nil {
		log.Println("failed to receive the message:" + err.Error())
	} else {
		log.Println("message received:" + id)
		var eventsList events.EventList
		err = proto.Unmarshal(msg.Content, &eventsList)
		if err != nil {
			log.Println("err:" + err.Error())
			return
		}
		eventStream <- &eventsList
	}
}

func NewSubscriber(validatorUrl string, db *DbServer) *Subscriber {
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
		done:       make(chan interface{}),
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
	msgType, msg, err := s.checkResponse(corId)
	if msgType == validator.Message_CLIENT_EVENTS_SUBSCRIBE_RESPONSE {
		res, ok := msg.(*client.ClientEventsUnsubscribeResponse)
		if !ok || res.Status != client.ClientEventsUnsubscribeResponse_OK {
			log.Println("failed to parse out the reponse")
		}
	} else {
		log.Fatal("corrupted response from:Message_CLIENT_EVENTS_UNSUBSCRIBE_REQUEST")
	}
}
