package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	server "github.com/rico-bee/leopark/market_api/server"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	rpcUrl = "localhost:50051"
)

var (
	version string
)

func main() {

	showversion := kingpin.Flag("version", "Show version information.").Short('v').Bool()
	// configPath := kingpin.Flag("config", "Optional config file path.").Default("config.user.json").String()
	debug := kingpin.Flag("debug", "Set logging to debug.").Bool()
	kingpin.Parse()

	if *showversion == true {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(rpcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMarketClient(conn)

	apiServer, err := server.NewServer(version, *debug)
	if err != nil {
		fmt.Printf("Failed to create Server: %s\n", err.Error())
		os.Exit(1)
	}
	apiServer.Start()
}
