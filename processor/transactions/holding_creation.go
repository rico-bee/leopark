package transactions

import (
	"errors"
	"fmt"
	pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	pb "github.com/rico-bee/leopark/market"
	"log"
)

func handleHoldingCreation(createHolding *pb.CreateHolding, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	acc, err := state.GetAccount(header.SignerPublicKey)
	if err != nil {
		log.Println("cannot get account")
	}
	if acc == nil {
		accErr := fmt.Sprintf("Account with public key %s doesn't exists", header.SignerPublicKey)
		return nil, errors.New(accErr)
	}
	asset := state.GetAsset(createHolding.Asset)
	if asset == nil {
		assetErr := fmt.Sprintf("Asset with name %s doesn't exists", createHolding.Asset)
		return nil, errors.New(assetErr)
	}

	state.CreateHolding(createHolding.Id, createHolding.Label, createHolding.Description,
		header.SignerPublicKey, createHolding.Asset, createHolding.Quantity)
	return state.UpdateHolding(header.SignerPublicKey, createHolding.Quantity)
}
