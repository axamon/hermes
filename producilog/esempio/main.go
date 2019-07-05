package main

import (
	"context"
	"log"
	"os"

	"github.com/axamon/hermes/producilog"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logfile := os.Args[1]

	err := producilog.KafkaLocalProducer(ctx, logfile)
	if err != nil {
		log.Printf("ERROR impossibile produrre messaggi in kafka: %s\n", err.Error())
	}

}
