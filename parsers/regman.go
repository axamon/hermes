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

package parsers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

const headerregman = "#giornoq,cpeid,tgu,trap_timestamp,deviceid,devicetype,mode,originipaddress,averagebitrate,avgsskbps,bufferingduration,callerclass,callerrorcode,callerrormessage,callerrortype,callurl,errordesc,errorreason,eventname,levelbitrates,linespeedkbps,maxsschunkkbps,maxsskbps,minsskbps,streamingtype,videoduration,videoposition,videotitle,videotype,videourl,eventtype,fwversion,networktype,ra_version,update_time,trap_provider,mid,service_id,service_id_version,date_rif,video_provider,max_upstream_net_latency,min_upstream_net_latency,avg_upstream_net_latency,max_downstream_net_latency,min_downstream_net_latency,avg_downstream_net_latency,max_platform_latency,min_platform_latency,avg_platform_latency,packet_loss,preloaded_app_v"

var loc, _ = time.LoadLocation("Europe/Rome")

const timeRegmanFormat = "2006-01-02 15:04:05"

var isREGMAN = regexp.MustCompile(`(?m)^.*deviceid.*$`)

var regmanLock sync.Mutex

var regmanrecords []string

// REGMAN è il parser delle trap provenienti da REGMAN.
func REGMAN(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Apre nuovo file per salvare dati elaborati.
	newFile := strings.Split(logfile, ".csv.gz")[0] + ".offuscato.csv.gz"
	// fmt.Println(newFile)

	f, err := os.Create(newFile)
	if err != nil {
		log.Println(err.Error())
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	// Scrive headers.
	gw.Write([]byte("#Log REGMAN prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte(headerregman + "\n"))

	// fmt.Println(logfile) // debug

	// Apri file zippato in memoria.
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file REGMAN %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	var records, s []string
	// var topic string

	n := 0
	for scan.Scan() {
		n++

		// Salta header.
		if n == 1 {
			continue
		}

		line := scan.Text()

		// Verifica che logfile sia di tipo regman.
		// if !isREGMAN.MatchString(line) {
		// 	err := fmt.Errorf("Error logfile %s non di tipo REGMAN", logfile)
		// 	return err
		// }

		_, s, err = elaboraREGMAN(ctx, &line)
		if err != nil {
			log.Printf("Error Impossibile elaborare REGMAN record: %s", s)
		}

		// if len(s) < 2 {
		// 	continue
		// }

		//fmt.Println(s[:]) // debug

		// Viene aggiunto a records un record con i campi individuati
		// separati da ",".

		// Scrive dati.
		justString := strings.Join(s, ",")
		// fmt.Println(justString)
		gw.Write([]byte(justString + "\n"))

	}

	// Scrive footer.
	gw.Write([]byte("#Numero di records: " + strconv.Itoa(len(records)) + "\n"))
	gw.Close()

	// Scrive uno per uno su standard output i record offuscati.
	// for _, record := range records {
	// 	fmt.Println(record)
	// }

	// Invia i records su kafka locale.
	// err = inoltralog.LocalKafkaProducer(ctx, topic, records)
	// if err != nil {
	// 	log.Printf("Error Impossibile salvare su kafka: %s\n", err.Error())
	// }

	fmt.Println(n)
	return err
}

func elaboraREGMAN(ctx context.Context, line *string) (topic string, result []string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// ricerca le fruzioni nell'intervallo temporale richiesto
	// l'intervallo temporale inzia con l'inzio di una fruizione

	// Splitta la linea nei supi fields.
	// Il separatore per i log REGMAN è ";"
	s := strings.Split(*line, ";")

	t, err := time.ParseInLocation(timeRegmanFormat, s[2], loc)
	if err != nil {
		log.Println(err.Error())
	}

	ora := t.UTC().Hour()
	minuto := t.UTC().Minute()

	// calcola a quale quartodora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15

	quartooraStr := strconv.Itoa(quartoora)

	//IDipq, _ := hasher.StringSum(s[6] + quartooraStr)

	//epoch := t.Format(time.RFC1123Z)

	// Crea il campo giornoq per integrare i log al quarto d'ora.
	giornoq := t.UTC().Format("20060102") + "q" + quartooraStr

	//Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)

	// recupera ip cliente

	//! OFFUSCAMENTO CAMPI SENSIBILI
	// s[6] contiente ip pubblico cliente.
	s[6], err = hasher.StringSumWithSalt(s[6], salt)
	if err != nil {
		log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
	}

	// s[1] contiene il cli del cliente.
	s[1], err = hasher.StringSumWithSalt(s[1], salt)
	if err != nil {
		log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
	}

	e := strings.Join(s, ",")

	result = append(result, giornoq, e)

	return giornoq, result, err
}
