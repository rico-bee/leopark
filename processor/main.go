package main

import (
	processorSdk "github.com/hyperledger/sawtooth-sdk-go/processor"
	transactions "github.com/rico-bee/leopark/processor/transactions"
	"log"
)

func main() {
	processor := processorSdk.NewTransactionProcessor("tcp://localhost:4040") // pass in validator url
	handler := transactions.MarketplaceHandler{}
	processor.AddHandler(&handler)
	err := processor.Start()
	if err != nil {
		log.Fatal("Failed to start processor")
	}
}
