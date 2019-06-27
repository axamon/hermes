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
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

var isCDN = regexp.MustCompile(`(?s)^\[.*\]\t[0-9]+\t\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\t[A-Z_]+\/\d{3}\t\d+\t[A-Z]+\t.*$`)

var chanRecords = make(chan *[]string)

var n int

var wg sync.WaitGroup

// CDN è il parser dei log provenienti dalla Content Delivery Network
func CDN(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var done = make(chan bool)
	// fmt.Println(logfile) // debug

	// Apri file zippato in memoria

	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		return
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	var records []string
	//var topic string

	go func() {
		select {
		case record := <-chanRecords:
			//fmt.Println(s[:]) // debug
			s := *record
			if len(s) < 2 {
				break
			}
			// ! OFFUSCAMENTO IP PUBBLICO CLIENTE
			// s[2] è l'ip pubblico del cliente da offuscare
			s[2], err = hasher.StringSumWithSalt(s[2], salt)
			if err != nil {
				log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
			}
			records = append(records, strings.Join(s, ","))
		case <-done:
			break
		}
		return
	}()

	for scan.Scan() {
		n++
		fmt.Println(n)
		line := scan.Text()

		// Verifica che logfile sia di tipo CDN.
		// if !isCDN.MatchString(line) {
		// 	err := fmt.Errorf("Error logfile %s non di tipo CDN", logfile)
		// 	return err
		// }
		wg.Add(1)

		go elaboraCDN(ctx, &line)

		// if err != nil {
		// 	log.Printf("Error Impossibile elaborare fruzione per record: %s", s)
		// }

	}

	wg.Wait()
	done <- true
	// Apre nuovo file per salvare dati elaborati.
	newFile := strings.Split(logfile, ".csv.gz")[0] + ".offuscato.csv.gz"
	// fmt.Println(newFile)

	f, err := os.Create(newFile)
	if err != nil {
		log.Println(err.Error())
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	justString := strings.Join(records, "\n")
	// fmt.Println(justString)

	// Scrive headers.
	gw.Write([]byte("#Log CDN prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte("#giornoq,hashfruizione,clientip,idvideoteca,status,tts[nanosecondi],bytes[bytes]\n"))
	// Scrive dati.
	gw.Write([]byte(justString + "\n"))
	// Scrive footer.
	gw.Write([]byte("#Numero di records: " + strconv.Itoa(len(records)) + "\n"))
	gw.Close()

	// Scrive uno per uno su standard output i record offuscati.
	// for _, line := range records {
	// 	fmt.Println(line)
	// }

	// Invia i records su kafka locale.
	//err = inoltralog.LocalKafkaProducer(ctx, topic, records)
	// err = inoltralog.RemoteKafkaProducer(ctx, "52.157.136.139:9092", topic, records)
	// if err != nil {
	// 	log.Printf("Error Impossibile salvare su kafka: %s\n", err.Error())
	// }

	// fmt.Println(n)
	return err
}

func elaboraCDN(ctx context.Context, line *string) { //(topic string, result []string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer wg.Done()
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

	ora := t.Hour()
	minuto := t.Minute()

	// Calcola a quale quarto d'ora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15

	// Trasforma quartoora in stringa.
	quartooraStr := strconv.Itoa(quartoora)

	//IDipq, _ := hasher.StringSum(s[2] + quartooraStr)

	//epoch := t.Format(time.RFC1123Z)

	//Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)
	//var speed, tts, bytes float64

	// Crea il campo giornoq per integrare i log al quarto d'ora.
	giornoq := t.Format("20060102") + "q" + quartooraStr

	// tts, err = strconv.ParseFloat(s[1], 8)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// bytes, err = strconv.ParseFloat(s[4], 8)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// Calcola la velocità di download.
	//speed = (bytes / tts)

	// Trasforma la velocità in stringa.
	//speedStr := fmt.Sprintf("%f", speed)

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
	var result []string
	result = append(result, giornoq, Hashfruizione, clientip, idvideoteca, status, s[1], s[4])
	chanRecords <- &result
	return
}
