package transactions

import (
	"errors"
	"fmt"
	pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	pb "github.com/rico-bee/leopark/market"
	"log"
)

func handleAssetCreation(createAsset *pb.CreateAsset, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	acc, err := state.GetAccount(header.SignerPublicKey)
	if err != nil {
		log.Println("cannot find account")
	}
	if acc == nil {
		msg := fmt.Sprintf("Account with public key %s doesn't exists", header.SignerPublicKey)
		return []string{}, errors.New(msg)
	}
	asset := state.GetAsset(createAsset.Name)
	if asset != nil {
		msg := fmt.Sprintf("Asset with name %s already exists", createAsset.Name)
		return []string{}, errors.New(msg)
	}
	return state.SetAsset(createAsset.Name, createAsset.Description, []string{header.SignerPublicKey}, createAsset.Rules)
}
