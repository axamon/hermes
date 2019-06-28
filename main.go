package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/parsers"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Così evitiamo problemi con l'istanzioamento degli errori :)
	var err error

	// // Verifica se l'instaza locale di Kakfa è raggiungile.
	// err = inoltralog.VerificaLocalKafka(ctx)
	// if err != nil {
	// 	log.Printf("ERROR Attenzione istanza locale di Kafka non raggiugibile\n")
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println("Ok proseguo lo stesso!")
	// }

	if len(os.Args) < 3 {
		fmt.Println("Devi inserire il filename e il tipo di log [CDN AVS REGMAN]")
		os.Exit(1)
	}

	logfile := os.Args[1]
	tipo := os.Args[2]

	switch tipo {
	case "CDN":
		err = parsers.CDN(ctx, logfile)
		if err != nil {
			log.Printf("Error Impossibile parsare file CDN %s: %s\n", logfile, err.Error())
		}
	case "REGMAN":
		err = parsers.REGMAN(ctx, logfile)
		if err != nil {
			log.Printf("Error Impossibile parsare file REGMAN %s: %s\n", logfile, err.Error())
		}
	case "AVS":
		err = parsers.AVS(ctx, logfile)
		if err != nil {
			log.Printf("Error Impossibile parsare file AVS %s: %s\n", logfile, err.Error())
		}
	default:
		fmt.Println("Specifica tipo di file: [CDN AVS REGMAN]")
	}
}
