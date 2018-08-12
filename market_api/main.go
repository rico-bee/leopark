package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	api "github.com/rico-bee/leopark/market_api/api"
	server "github.com/rico-bee/leopark/market_api/server"
	pb "github.com/rico-bee/leopark/market_service/proto/api"
	"google.golang.org/grpc"
	"log"
	"os"
)

const (
	version = "0.0.1"
	rpcUrl  = "localhost:50051"
)

func main() {
	showversion := kingpin.Flag("version", "Show version information.").Short('v').Bool()
	if *showversion == true {
		log.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	serviceUrl := kingpin.Flag("service", "Service Url.").Default(rpcUrl).Short('s').String()
	if serviceUrl == nil {
		log.Fatal("no service url is defined")
	}

	rethinkdbUrl := kingpin.Flag("database", "rethinkdb Url.").Default("localhost:28015").Short('r').String()
	if rethinkdbUrl == nil {
		log.Fatal("no database url is defined")
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*serviceUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMarketClient(conn)

	// init db connection
	db, err := api.NewDBServer(*rethinkdbUrl)
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
