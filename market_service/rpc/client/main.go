package main

import (
	"log"
	"time"

	pb "github.com/rico-bee/marketplace/market_service/proto/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	url   = "localhost:50051"
	name  = "rico"
	email = "ricozhang726"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMarketClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.DoCreateAccount(ctx, &pb.CreateAccountRequest{Name: name, Email: email})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf("Create user: %s", r.Message)
}
