package transactions

import (
	"errors"
	"fmt"
	pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	pb "github.com/rico-bee/leopark/market"
	"log"
)

func handleAccountCreation(createAcc *pb.CreateAccount, header *pb2.TransactionHeader, state *MarketState) ([]string, error) {
	acc, err := state.GetAccount(header.SignerPublicKey)
	if err != nil {
		log.Println("warning:cannot find the account")
	}
	if acc != nil {
		msg := fmt.Sprintf("Account with public key %s already exists", header.SignerPublicKey)
		return nil, errors.New(msg)
	}
	return state.SetAccount(header.SignerPublicKey, createAcc.Label, createAcc.Description, []string{})
}
