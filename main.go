// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
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
