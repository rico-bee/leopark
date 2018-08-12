package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"os/signal"
)

const (
	version = "0.0.1"
)

func main() {
	app := kingpin.New("market sync", "A command-line sync application sync data to db.")
	showversion := app.Flag("version", "Show version information.").Short('v').Bool()
	if *showversion == true {
		log.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	validatorUrl := app.Flag("validator", "Validator Url.").Default("tcp://localhost:4040").Short('d').String()
	if validatorUrl == nil {
		log.Println("no validator url is defined")
		*validatorUrl = "tcp://localhost:4040"
	}

	rethinkdbUrl := app.Flag("database", "rethinkdb Url.").Default("localhost:28015").Short('r').String()
	if rethinkdbUrl == nil {
		log.Println("no validator url is defined")
		*rethinkdbUrl = "localhost:28015"
	}
	app.Parse(os.Args[1:])
	log.Printf("db: %s\n", *rethinkdbUrl)
	log.Printf("validator: %s\n", *validatorUrl)
	db, err := NewDBServer(*rethinkdbUrl)
	if err != nil {
		log.Fatal("failed to connect to DB")
	}
	knownBlocks, err := db.LastKnownBlocks(15)
	if err != nil {
		log.Println("failed to find last known blocks")
	}

	done := make(chan interface{})
	subscriber := NewSubscriber(*validatorUrl, done, db)
	subscriber.Start(knownBlocks)

	c := make(chan os.Signal, 1)
	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	signal.Notify(c)

	s := <-c
	close(done)
	log.Println("gracefully exit:" + s.String())
}
