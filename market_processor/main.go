package main

import (
	processorSdk "github.com/hyperledger/sawtooth-sdk-go/processor"
	transactions "github.com/rico-bee/leopark/market_processor/transactions"
	"log"
	"os"
	// "os/signal"
	// "syscall"
)

func main() {
	processor := processorSdk.NewTransactionProcessor("tcp://localhost:4040") // pass in validator url
	handler := transactions.MarketplaceHandler{}
	processor.AddHandler(&handler)
	processor.ShutdownOnSignal(os.Interrupt, os.Kill)
	err := processor.Start()
	if err != nil {
		log.Fatal("cannot start processor:" + err.Error())
	}
}
