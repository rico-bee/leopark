package main

import (
	"github.com/alecthomas/kingpin"
	processorSdk "github.com/hyperledger/sawtooth-sdk-go/processor"
	transactions "github.com/rico-bee/leopark/market_processor/transactions"
	"log"
	"os"
)

const (
	version = "0.0.1"
)

func main() {
	showversion := kingpin.Flag("version", "Show version information.").Short('v').Bool()
	if *showversion == true {
		log.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	validatorUrl := kingpin.Flag("validator", "Validator Url.").Default("tcp://localhost:4040").Short('d').String()
	if validatorUrl == nil {
		log.Fatal("no validator url is defined")
	}
	kingpin.Parse()
	processor := processorSdk.NewTransactionProcessor(*validatorUrl) // pass in validator url
	handler := transactions.MarketplaceHandler{}
	processor.AddHandler(&handler)
	processor.ShutdownOnSignal(os.Interrupt, os.Kill)
	err := processor.Start()
	if err != nil {
		log.Fatal("cannot start processor:" + err.Error())
	}
}
