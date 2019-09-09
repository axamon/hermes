// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/axamon/hermes/parsers"
)

var stat = flag.Bool("stat", false, "riporta a video delle statistiche")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var goroutines = flag.Int("g", 1000, "Numero di processi paralleli da usare")
var logfile = flag.String("f", "", "Logfile da parsare")
var tipo = flag.String("t", "", "tipo Logfile da parsare")

func main() {

	flag.Parse()

	var start = time.Now()
	

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

	// Così evitiamo problemi con l'istanzioamento degli errori :)
	var err error

	switch *tipo {
	case "CDN":
		err = parsers.CDN(ctx, *logfile, *goroutines)
		if err != nil {
			log.Printf("Error Impossibile parsare file CDN %s: %s\n", *logfile, err.Error())
		}
	case "REGMAN":
		err = parsers.REGMAN(ctx, *logfile, *goroutines)
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

	if *stat {
		// start := time.Now()
		log.Printf("Elaborazione %s: %v\n", *logfile, time.Since(start))
	}
}
