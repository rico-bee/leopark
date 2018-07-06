package transaction

import (
	addresser "github.com/rico-bee/marketplace/address"
	pb "github.com/rico-bee/marketplace/market"
	"crypto/sha512"
	"encoding/hex"
	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"log"
)

// TxHeader : Transaction Header in
type TxHeader = transaction_pb2.TransactionHeader
type Signer = signing.Signer
type BatchHeader = batch_pb2.BatchHeader
type Batch = batch_pb2.Batch

func makeHeader(inputAddresses []string, outputAddresses []string, payloadSha512 string,
	signerPublicKey string, batcherPublicKey string) *TxHeader {
	nonce, err := uuid.NewRandom() // move to a lib
	if err != nil {
		log.Println("whatever it is wrong..")
	}

	header := &transaction_pb2.TransactionHeader{
		Inputs:           inputAddresses,
		Outputs:          outputAddresses,
		BatcherPublicKey: batcherPublicKey,
		Dependencies:     []string{},
		FamilyName:       addresser.FamilyName,
		FamilyVersion:    "1.0",
		Nonce:            nonce.String(),
		SignerPublicKey:  signerPublicKey,
		PayloadSha512:    payloadSha512,
	}

	return header
}

func makeHeaderAndBatch(payload []byte, inputAddresses []string, outputAddresses []string,
	txnKey *Signer, batchKey *Signer) ([]*batch_pb2.Batch, string) {

	payloadSha := sha512.Sum512(payload)
	payloadHexStr := hex.EncodeToString(payloadSha[:])
	header := makeHeader(inputAddresses, outputAddresses,
		payloadHexStr,
		txnKey.GetPublicKey().AsHex(),
		batchKey.GetPublicKey().AsHex())

	headerBytes, err := proto.Marshal(header)
	if err != nil {
		log.Fatal("Corrupted header: " + err.Error())
	}
	txnKeyStr := hex.EncodeToString(txnKey.Sign(headerBytes))
	transaction := &transaction_pb2.Transaction{
		Payload:         payload,
		Header:          headerBytes,
		HeaderSignature: txnKeyStr,
	}

	txnbytes, err := proto.Marshal(transaction)
	if err != nil {
		log.Println("txnbytes:" + string(txnbytes))
	}
	batchHeader := &batch_pb2.BatchHeader{
		SignerPublicKey: batchKey.GetPublicKey().AsHex(),
		TransactionIds:  []string{transaction.HeaderSignature},
	}

	batchHeaderBytes, err := proto.Marshal(batchHeader)
	if err != nil {
		log.Println("invalid batchHeader" + err.Error())
	}
	headerSignatureStr := hex.EncodeToString(batchKey.Sign(batchHeaderBytes))
	batch := &batch_pb2.Batch{
		Header:          batchHeaderBytes,
		HeaderSignature: headerSignatureStr,
		Transactions:    []*transaction_pb2.Transaction{transaction},
	}
	return []*batch_pb2.Batch{batch}, batch.HeaderSignature
}

func CreateAccount(txnKey *Signer, batchKey *Signer, label, description string) ([]*batch_pb2.Batch, string) {
	inputs := []string{addresser.MakeAccountAddress(txnKey.GetPublicKey().AsHex())}
	outputs := []string{addresser.MakeAccountAddress(txnKey.GetPublicKey().AsHex())}

	createAccount := &pb.CreateAccount{
		Label:       label,
		Description: description,
	}
	payload := &pb.TransactionPayload{
		PayloadType:   pb.TransactionPayload_CREATE_ACCOUNT,
		CreateAccount: createAccount,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal payload")
	}
	return makeHeaderAndBatch(data, inputs, outputs, txnKey, batchKey)
}
