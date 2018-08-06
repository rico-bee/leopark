package transaction

import (
	"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	batch_pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	batch_submit "github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_submit_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"

	zmq "github.com/pebbe/zmq4"
	//"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	userAgent              string = "rpc-server"
	httpMaxIdleConnections int    = 30
	httpRequestTimeout     int    = 120
	httpServerReadTimeout  int    = 61
	httpServerWriteTimeout int    = 120
)

type Payload []byte

type SawtoothAPI struct {
	connection messaging.Connection
}

func newTransport() *http.Transport {
	return &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        30,
	}
}

func createHTTPClient() *http.Client {
	return &http.Client{
		Transport: newTransport(),
		Timeout:   time.Duration(120) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (s *SawtoothAPI) checkResponse(corId string, message proto.Message) (validator_pb2.Message_MessageType, error) {
	_, msg, err := s.connection.RecvMsgWithId(corId)
	log.Println("received message:" + corId)
	if err != nil {
		log.Println(err.Error())
		return validator_pb2.Message_DEFAULT, err
	}

	err = proto.Unmarshal(msg.Content, message)
	if err != nil {
		return validator_pb2.Message_DEFAULT, err
	}
	return msg.MessageType, nil
}

func (s *SawtoothAPI) CheckBatchStatus(batchIds []string) (bool, error) {
	req := &batch_submit.ClientBatchStatusRequest{
		BatchIds: batchIds,
		Wait:     true,
		Timeout:  600,
	}
	rawBytes, err := proto.Marshal(req)
	if err != nil {
		log.Println("corrupted msg" + err.Error())
	}

	res := &batch_submit.ClientBatchStatusResponse{}
	id, err := s.connection.SendNewMsg(validator_pb2.Message_CLIENT_BATCH_STATUS_REQUEST, rawBytes)
	if err != nil {
		log.Println("failed to check status:" + err.Error())
	}
	_, err = s.checkResponse(id, res)
	if err != nil {
		log.Println("failed to query batch status")
	}
	for _, b := range res.BatchStatuses {
		log.Println("batch status:" + b.Status.String())
	}
	status := res.BatchStatuses[0].Status
	if status == batch_submit.ClientBatchStatus_COMMITTED {
		return true, nil
	}
	return false, errors.New("batch is not committed yet." + status.String())
}

func (s *SawtoothAPI) BatchRequest(batches []*batch_pb2.Batch) error {
	req := &batch_submit.ClientBatchSubmitRequest{
		Batches: batches,
	}
	rawBytes, err := proto.Marshal(req)
	if err != nil {
		log.Println("corrupted msg" + err.Error())
	}
	batchId, err := s.connection.SendNewMsg(validator_pb2.Message_CLIENT_BATCH_SUBMIT_REQUEST, rawBytes)
	if err != nil {
		log.Println("failed to send batch request to validator:" + err.Error())
	}
	log.Println("batch id:" + batchId)
	return err
}

func NewSawtoothApi(validatorUrl string) *SawtoothAPI {
	zmqCtx, err := zmq.NewContext()
	if err != nil {
		log.Fatalln("cannot create zmq context")
	}

	conn, err := messaging.NewConnection(zmqCtx, zmq.DEALER, validatorUrl, false)
	if err != nil {
		log.Println("failed to create connection")
	}
	return &SawtoothAPI{
		connection: conn,
	}
}
