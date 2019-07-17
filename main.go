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
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/axamon/hermes/parsers"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var logfile = flag.String("f", "", "Logfile da parsare")
var tipo = flag.String("t", "", "tipo Logfile da parsare")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parsers.AVS(ctx, "/home/gioml/Analytics/new/AVS/Fruizioni_08062019.csv.gz")

	// Così evitiamo problemi con l'istanzioamento degli errori :)
	var err error

	switch *tipo {
	case "CDN":
		err = parsers.CDN(ctx, *logfile)
		if err != nil {
			log.Printf("Error Impossibile parsare file CDN %s: %s\n", *logfile, err.Error())
		}
	case "REGMAN":
		err = parsers.REGMAN(ctx, *logfile)
		if err != nil {
			log.Printf("Error Impossibile parsare file REGMAN %s: %s\n", *logfile, err.Error())
		}
	case "AVS":
		err = parsers.AVS(ctx, *logfile)
		if err != nil {
			log.Printf("Error Impossibile parsare file AVS %s: %s\n", *logfile, err.Error())
		}
	default:
		fmt.Println("Specifica tipo di file: [CDN AVS REGMAN]")
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
