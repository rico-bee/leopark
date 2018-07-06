package transactions

import (
	pb "github.com/rico-bee/marketplace/market"
	"errors"
	"fmt"
	pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"log"
)

func handleAccountCreation(createAcc *pb.CreateAccount, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	acc, err := state.GetAccount(header.SignerPublicKey)
	if err != nil {
		//return nil, errors.New("cannot find the account")
		log.Println("cannot find the account")
	}
	if acc != nil {
		msg := fmt.Sprintf("Account with public key %s already exists", header.SignerPublicKey)
		return nil, errors.New(msg)
	}
	return state.SetAccount(header.SignerPublicKey, createAcc.Label, createAcc.Description, []string{})
}
