package main

import (
	rpc "bitbucket.org/riczha/marketplace/market_service/rpc"
	"fmt"
	"github.com/alecthomas/kingpin"
	"os"
)

var (
	version string
)

func main() {

	showversion := kingpin.Flag("version", "Show version information.").Short('v').Bool()
	configPath := kingpin.Flag("config", "Optional config file path.").Default("config.user.json").String()
	kingpin.Parse()

	if *showversion == true {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	_, err := os.Stat(*configPath)
	if err != nil {
		fmt.Printf("Config file: %s not found.\n", *configPath)
		os.Exit(1)
	}
	rpc.StartRpcServer()
}
