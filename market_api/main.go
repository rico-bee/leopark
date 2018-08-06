package main

import (
	"fmt"
	//"github.com/alecthomas/kingpin"
	server "github.com/rico-bee/leopark/market_api/server"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
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
	// Set up a connection to the server.
	conn, err := grpc.Dial(rpcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMarketClient(conn)
	// Contact the server and print out its response.
	apiServer, err := server.NewServer(c)

	if err != nil {
		fmt.Printf("Failed to create Server: %s\n", err.Error())
		os.Exit(1)
	}
	apiServer.Start()
}
