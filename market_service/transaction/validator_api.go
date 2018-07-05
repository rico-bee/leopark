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

// func (s *SawtoothAPI) CheckBatchRequest(batchId string) error {
// 	req := &client_batch_submit_pb2.ClientBatchStatusRequest{
// 		BatchIds: []string{batchId},
// 	}
// }

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

// func (s *SawtoothAPI) request(path, method string, payload interface{}) (Payload, error) {
// 	url := s.validatorURL + "/" + path
// 	buf := new(bytes.Buffer)
// 	var contentType string
// 	switch v := payload.(type) {
// 	case string:
// 		buf = bytes.NewBufferString(v)
// 		contentType = MIMEApplicationForm
// 	default:
// 		err := json.NewEncoder(buf).Encode(payload)
// 		if err != nil {
// 			return nil, err
// 		}
// 		contentType = MIMEApplicationJSON
// 	}
// 	req, _ := http.NewRequest(method, url, buf)
// 	req.Header.Add(headerAccept, MIMEApplicationJSON)
// 	req.Header.Add(headerContentType, contentType)
// 	req.Header.Add(headerUserAgent, userAgent)

// 	resp, err := s.httpClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	respContentType := resp.Header.Get(headerContentType)
// 	if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != 503 && respContentType == MIMEApplicationJSON {
// 		log.Fatalf("failed to CALL validator")
// 		return nil, errors.New("failed to call the validator")
// 	}
// 	// only 200 || 201 || 202 is expected from this call
// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
// 		return nil, errors.New("failed to call the validator")
// 	}
// 	bodyBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bodyBytes, nil
// }
