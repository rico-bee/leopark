package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	rpc "github.com/rico-bee/leopark/market_service/rpc"
	"os"
)

var (
	version string = "0.0.1"
)

func main() {

	showversion := kingpin.Flag("version", "Show version information.").Short('v').Bool()
	validatorUrl := kingpin.Flag("validator", "Validator Url.").Default("tcp://localhost:4040").Short('d').String()
	kingpin.Parse()

	if *showversion == true {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	if validatorUrl == nil {
		panic("no validator url is defined")
	}
	rpc.StartRpcServer(*validatorUrl)
}
