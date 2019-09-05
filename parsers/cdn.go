// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

var isCDN = regexp.MustCompile(`(?s)^\[.*\]\t[0-9]+\t\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\t[A-Z_]+\/\d{3}\t\d+\t[A-Z]+\t.*$`)

var n int

var cdnrecords []string
var cdnRecords sync.Mutex

// CDN è il parser dei log provenienti dalla Content Delivery Network
func CDN(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	start := time.Now()
	//	var done = make(chan bool)
	// fmt.Println(logfile) // debug

	// Apri file zippato in memoria

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
	//gw.Write([]byte("#Log CDN prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte("giornoq;hashfruizione;clientip;idvideoteca;status;tts[nanosecondi];bytes[bytes]\n"))

	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		return
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	//var topic string

	for scan.Scan() {
		n++
		line := scan.Text()

		// Verifica che logfile sia di tipo CDN.
		//if !isCDN.MatchString(line) {
		//	err := fmt.Errorf("Error logfile %s non di tipo CDN", logfile)
		//	return err
		//}

		elaboraCDN(ctx, &line, gw)

		// if err != nil {
		// 	log.Printf("Error Impossibile elaborare fruzione per record: %s", s)
		// }

	}
	//done <- true

	//justString := strings.Join(cdnrecords, "\n")
	// fmt.Println(justString)

	// Scrive dati.
	//gw.Write([]byte(justString + "\n"))
	// Scrive footer.
	// cdnRecords.Lock()
	// _, err = gw.Write([]byte("#Numero di cdnrecords: " + strconv.Itoa(n) + "\n"))
	// if err != nil {
	// 	log.Println("ERROR Impossibile scrivere su file zippato")
	// }
	// cdnRecords.Unlock()
	gw.Flush()
	gw.Close()

	// Scrive uno per uno su standard output i record offuscati.
	// for _, line := range cdnrecords {
	// 	fmt.Println(line)
	// }

	// Invia i cdnrecords su kafka locale.
	//err = inoltralog.LocalKafkaProducer(ctx, topic, cdnrecords)
	// err = inoltralog.RemoteKafkaProducer(ctx, "52.157.136.139:9092", topic, cdnrecords)
	// if err != nil {
	// 	log.Printf("Error Impossibile salvare su kafka: %s\n", err.Error())
	// }

	// fmt.Println(n)
	fmt.Println("Impiegato: ", time.Since(start))
	return err
}

func elaboraCDN(ctx context.Context, line *string, gw *gzip.Writer) { //(topic string, result []string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// ricerca le fruzioni nell'intervallo temporale richiesto
	// l'intervallo temporale inzia con l'inzio di una fruizione

	// fmt.Println(line)
	//if !isCDN.MatchString(line) {
	//err := fmt.Errorf("Error record non di tipo CDN: %s", line)
	//	return
	//}

	// Splitta la linea nei suoi fields,
	// il separatore per i log CDN è il tab: \t
	s := strings.Split(*line, "\t")

	// Parsa la URL nelle sue componenti.
	u, err := url.Parse(s[6])
	if err != nil {
		log.Printf("Error nel parsing URL di: %s\n", *line)
	}
	/* Urlschema := u.Scheme
	if Urlschema != "https" { //fa passare solo le URL richieste via WEB
		continue
	} */

	// Converte il timestamp del log.
	t, err := time.Parse("[02/Jan/2006:15:04:05.000+000]", s[0])
	if err != nil {
		log.Println(err.Error())
	}

	giornoq := giornoq(t)

	// Recupera l'ip del cliente.
	clientip := s[2]

	// Recupera lo status HTTP del chunk.
	status := s[3]

	// Recupera lo user agent del cliente.
	ua := s[8]

	//fmt.Println(Urlschema)
	//Urlhost := u.Host
	Urlpath := u.Path
	//fmt.Println(Urlpath)
	//Urlquery := u.RawQuery
	//Urlfragment := u.Fragment
	pezziurl := strings.Split(Urlpath, "/")
	//fmt.Println(pezziurl)

	// Tratta solo i chunck di tipo video // ! da verificare se va bene o no!
	// if ok := !strings.Contains(Urlpath, "video="); ok == true { //solo i chunk video

	// 	return "", nil, nil
	// }
	if len(pezziurl) < 6 {
		return
	}
	// Recupera il valore univoco del video.
	idvideoteca := pezziurl[6]

	//tipocodifica := pezziurl[7]
	//idavs := pezziurl[8]
	//fmt.Println(idvideoteca)
	//encoding := pezziurl[10]
	//fmt.Println(encoding)
	//re := regexp.MustCompile(`QualityLevels\(([0-9]+)\)$`)
	//bitratestr := re.FindStringSubmatch(encoding)[1]
	//bitrate, _ := strconv.ParseFloat(bitratestr, 8)
	/* if err != nil {
		log.Fatal(err.Error())
	} */
	//bitrateMB := bitrate * bitstoMB

	// Crea l'idfruzione univoco del cliente.
	Hashfruizione, err := hasher.StringSum(clientip + idvideoteca + ua)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	//ingestafruizioni(Hash, clientip, idvideoteca, idaps, edgeip, giorno, orario, speed)

	//s = append(s, Time, Hashfruizione, idaps, idvideoteca, status, speedStr, quartooraStr, IDipq)
	var str []string
	str = append(str, giornoq, Hashfruizione, clientip, idvideoteca, status, s[1], s[4])

	if len(str) < 2 {
		return
	}
	// ! OFFUSCAMENTO IP PUBBLICO CLIENTE
	// s[2] è l'ip pubblico del cliente da offuscare
	str[2], err = hasher.StringSumWithSalt(str[2], salt)
	if err != nil {
		log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
	}

	record := strings.Join(str, ";") + "\n"
	//cdnrecords = append(cdnrecords, strings.Join(str, ";"))
	cdnRecords.Lock()
		gw.Write([]byte(record))
	cdnRecords.Unlock()

	//chancdnrecords <- &result
	return
}
