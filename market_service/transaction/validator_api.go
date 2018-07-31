package transaction

import (
	//"errors"
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	batch_pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_submit_pb2"
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

func (s *SawtoothAPI) BatchRequest(batches []*batch_pb2.Batch) error {
	req := &client_batch_submit_pb2.ClientBatchSubmitRequest{
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
