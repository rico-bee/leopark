package main

import (
	"log"
	"os"
	"os/signal"
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

	done := make(chan interface{})
	subscriber := NewSubscriber("tcp://localhost:4040", done, db)
	subscriber.Start(knownBlocks)

	c := make(chan os.Signal, 1)
	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)

	s := <-c
	close(done)
	log.Println("gracefully exit:" + s.String())
}
