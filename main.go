package main

import (
	"log"
	"os"

	"github.com/axamon/hermes/parsers"
)

func main() {
	var err error
	logfile := os.Args[1]

	// fmt.Println(logfile)
	// err = parsers.CDN(logfile)
	// if err != nil {
	// 	log.Printf("Error Impossibile parsare file CDN %s: %s\n", logfile, err.Error())
	// }

	// err = parsers.REGMAN(logfile)
	// if err != nil {
	// 	log.Printf("Error Impossibile parsare file REGMAN %s: %s\n", logfile, err.Error())
	// }

	err = parsers.AVS(logfile)
	if err != nil {
		log.Printf("Error Impossibile parsare file REGMAN %s: %s\n", logfile, err.Error())
	}
}
