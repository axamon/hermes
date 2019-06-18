package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/axamon/hermes/inoltralog"
	"github.com/axamon/hermes/parsers"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Così evitiamo problemi con l'istanzioamento degli errori :)
	var err error

	// Verifica se l'instaza locale di Kakfa è raggiungile.
	err = inoltralog.VerificaLocalKafka(ctx)
	if err != nil {
		log.Printf("ERROR Attenzione istanza locale di Kafka non raggiugibile\n")
		time.Sleep(5 * time.Second)
		fmt.Println("Ok proseguo lo stesso!")
	}

	logfile := os.Args[1]

	fmt.Println(logfile)
	err = parsers.CDN(ctx, logfile)
	if err != nil {
		log.Printf("Error Impossibile parsare file CDN %s: %s\n", logfile, err.Error())
	}

	err = parsers.REGMAN(ctx, logfile)
	if err != nil {
		log.Printf("Error Impossibile parsare file REGMAN %s: %s\n", logfile, err.Error())
	}

	err = parsers.AVS(ctx, logfile)
	if err != nil {
		log.Printf("Error Impossibile parsare file REGMAN %s: %s\n", logfile, err.Error())
	}

}
