package transactions

import (
	"errors"
	"fmt"
	pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	pb "github.com/rico-bee/leopark/market"
	"log"
)

func handleAssetCreation(createAsset *pb.CreateAsset, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	log.Println("creating asset for account key : " + header.SignerPublicKey)
	acc, err := state.GetAccount(header.SignerPublicKey)

	if err != nil {
		log.Println("asset: cannot find account with key " + header.SignerPublicKey + " due to " + err.Error())
	}
	if acc == nil {
		msg := fmt.Sprintf("Account with public key %s doesn't exists", header.SignerPublicKey)
		log.Println(msg)
		return []string{}, errors.New(msg)
	}
	asset := state.GetAsset(createAsset.Name)
	if asset != nil {
		msg := fmt.Sprintf("Asset with name %s already exists", createAsset.Name)
		return []string{}, errors.New(msg)
	}
	log.Println("creating asset:" + createAsset.Name + ": " + createAsset.Description)
	return state.SetAsset(createAsset.Name, createAsset.Description, []string{header.SignerPublicKey}, createAsset.Rules)
}
