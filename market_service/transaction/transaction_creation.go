package transaction

import (
	"crypto/sha512"
	"encoding/hex"
	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	addresser "github.com/rico-bee/leopark/address"
	pb "github.com/rico-bee/leopark/market"
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

func CreateAsset(txnKey, batchKey *Signer, name, description string, rules []*pb.Rule) ([]*batch_pb2.Batch, string) {
	inputs := []string{addresser.MakeAssetAddress(txnKey.GetPublicKey().AsHex())}
	outputs := []string{addresser.MakeAssetAddress(txnKey.GetPublicKey().AsHex())}

	createAsset := &pb.CreateAsset{
		Name:        name,
		Description: description,
		Rules:       rules,
	}

	payload := &pb.TransactionPayload{
		PayloadType: pb.TransactionPayload_CREATE_ASSET,
		CreateAsset: createAsset,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal payload")
	}
	return makeHeaderAndBatch(data, inputs, outputs, txnKey, batchKey)
}

func CreateHolding(txnKey, batchKey *Signer, identifier, label, description, asset string, quantity int64) ([]*batch_pb2.Batch, string) {
	inputs := []string{addresser.MakeAccountAddress(txnKey.GetPublicKey().AsHex()), addresser.MakeAssetAddress(asset), addresser.MakeHoldingAddress(identifier)}
	outputs := []string{addresser.MakeAccountAddress(txnKey.GetPublicKey().AsHex()), addresser.MakeHoldingAddress(txnKey.GetPublicKey().AsHex())}

	createHolding := &pb.CreateHolding{
		Id:          identifier,
		Label:       label,
		Description: description,
		Asset:       asset,
		Quantity:    quantity,
	}
	payload := &pb.TransactionPayload{
		PayloadType:   pb.TransactionPayload_CREATE_HOLDING,
		CreateHolding: createHolding,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal payload")
	}
	return makeHeaderAndBatch(data, inputs, outputs, txnKey, batchKey)
}

type MarketplaceHolding struct {
	HoldingId string `json:"holdingId,omitempty"`
	Quantity  int64  `json:"quantity,omitempty"`
	Asset     string `json:"asset,omitempty"`
}

func CreateOffer(txnKey, batchKey *Signer, identifier, label, description string, source,
	target *MarketplaceHolding, rules []*pb.Rule) ([]*batch_pb2.Batch, string) {
	inputs := []string{addresser.MakeAccountAddress(txnKey.GetPublicKey().AsHex()),
		addresser.MakeAssetAddress(source.Asset), addresser.MakeOfferAddress(identifier)}
	outputs := []string{addresser.MakeOfferAddress(identifier),
		addresser.MakeHoldingAddress(txnKey.GetPublicKey().AsHex())}
	if target.HoldingId != "" {
		inputs = append(inputs, addresser.MakeHoldingAddress(target.HoldingId))
		inputs = append(inputs, addresser.MakeAssetAddress(target.Asset))
	}

	createOffer := &pb.CreateOffer{
		Id:             identifier,
		Label:          label,
		Description:    description,
		Source:         source.HoldingId,
		SourceQuantity: source.Quantity,
		Target:         target.HoldingId,
		TargetQuantity: target.Quantity,
		Rules:          rules,
	}

	payload := &pb.TransactionPayload{
		PayloadType: pb.TransactionPayload_CREATE_OFFER,
		CreateOffer: createOffer,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal payload")
	}
	return makeHeaderAndBatch(data, inputs, outputs, txnKey, batchKey)
}

type OfferParticipant struct {
	SrcHolding    string
	TargetHolding string
	SrcAsset      string
	TargetAsset   string
}

func AcceptOffer(txnKey, batchKey *Signer, identifier string, count uint64,
	sender, receiver *OfferParticipant) ([]*batch_pb2.Batch, string) {
	inputs := []string{
		addresser.MakeHoldingAddress(receiver.TargetHolding),
		addresser.MakeHoldingAddress(sender.SrcHolding),
		addresser.MakeAssetAddress(sender.SrcAsset),
		addresser.MakeAssetAddress(receiver.TargetAsset),
		addresser.MakeOfferHistoryAddress(identifier),
		addresser.MakeOfferAccountAddress(identifier, txnKey.GetPublicKey().AsHex()),
		addresser.MakeOfferAddress(identifier),
	}
	outputs := []string{
		addresser.MakeHoldingAddress(receiver.TargetHolding),
		addresser.MakeHoldingAddress(sender.SrcHolding),
		addresser.MakeOfferHistoryAddress(identifier),
		addresser.MakeOfferAccountAddress(identifier, txnKey.GetPublicKey().AsHex()),
	}

	if receiver.SrcHolding != "" {
		inputs = append(inputs, addresser.MakeHoldingAddress(receiver.SrcHolding))
		inputs = append(inputs, addresser.MakeAssetAddress(receiver.SrcAsset))
		outputs = append(outputs, addresser.MakeHoldingAddress(receiver.SrcHolding))
	}

	if sender.TargetHolding != "" {
		inputs = append(inputs, addresser.MakeHoldingAddress(sender.TargetHolding))
		inputs = append(inputs, addresser.MakeAssetAddress(sender.TargetAsset))
		outputs = append(outputs, addresser.MakeHoldingAddress(sender.TargetHolding))
	}
	acceptOffer := &pb.AcceptOffer{
		Id:     identifier,
		Source: receiver.SrcHolding,
		Target: receiver.TargetHolding,
		Count:  count,
	}

	payload := &pb.TransactionPayload{
		PayloadType: pb.TransactionPayload_ACCEPT_OFFER,
		AcceptOffer: acceptOffer,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal payload")
	}
	return makeHeaderAndBatch(data, inputs, outputs, txnKey, batchKey)
}

func CloseOffer(txnKey, batchKey *Signer, identifier string) ([]*batch_pb2.Batch, string) {
	inputs := []string{addresser.MakeOfferAddress(identifier)}

	outputs := []string{addresser.MakeOfferAddress(identifier)}

	closeOffer := &pb.CloseOffer{
		Id: identifier,
	}
	payload := &pb.TransactionPayload{
		PayloadType: pb.TransactionPayload_CLOSE_OFFER,
		CloseOffer:  closeOffer,
	}
	data, err := proto.Marshal(payload)
	if err != nil {
		log.Fatal("failed to marshal payload")
	}
	return makeHeaderAndBatch(data, inputs, outputs, txnKey, batchKey)
}
