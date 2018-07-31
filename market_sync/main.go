package main

import (
	"log"
)

func main() {
	db, err := NewDBServer("localhost:28015")
	if err != nil {
		log.Fatal("failed to connect to DB")
	}
	knownBlocks, err := db.LastKnownBlocks(15)
	if err != nil {
		log.Println("failed to find last known blocks")
	}
	subscriber := NewSubscriber("tcp://localhost:4040", db)
	subscriber.Start(knownBlocks)
}
