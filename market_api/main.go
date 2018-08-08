package main

import (
	"fmt"
	//"github.com/alecthomas/kingpin"
	api "github.com/rico-bee/leopark/market_api/api"
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

	// init db connection
	db, err := api.NewDBServer("localhost:28015")
	if err != nil {
		log.Fatal("failed to connect to DB")
	}

	// Contact the server and print out its response.
	apiServer, err := server.NewServer(&api.Handler{RpcClient: c, Db: db})

	if err != nil {
		fmt.Printf("Failed to create Server: %s\n", err.Error())
		os.Exit(1)
	}
	apiServer.Start()
}
