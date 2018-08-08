package main

import (
	processorSdk "github.com/hyperledger/sawtooth-sdk-go/processor"
	transactions "github.com/rico-bee/leopark/market_processor/transactions"
	"log"
	"os"
	"os/signal"
)

func main() {
	processor := processorSdk.NewTransactionProcessor("tcp://localhost:4040") // pass in validator url
	handler := transactions.MarketplaceHandler{}
	processor.AddHandler(&handler)
	err := processor.Start()
	c := make(chan os.Signal, 1)
	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)
	s := <-c
	processor.Shutdown()
	log.Println("gracefully exit:" + s.String())
	if err != nil {
		log.Fatal("Failed to start processor")
	}
}
