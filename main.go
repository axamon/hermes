package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/parsers"
)

func main() {

	ctx := context.WithCancel(context.Background(), cancel)
	defer cancel()

	var err error
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
