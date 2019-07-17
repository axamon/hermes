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
	"encoding/csv"
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

const avsheader = "#giornoq;device;timestamp;cli;hash;attivita;idvideoteca;standard;metodopagamento;ignoto;ignoto1;mailcliente;ignoto3;servizio;costruttore;tipoprog;case;rete;num"
const timeAVSFormat = "2006-01-02T15:04:05"

var isAVS = regexp.MustCompile(`(?m)^.*\|.*\|.*$`)

var avsRecords sync.Mutex

// AVS è il parser dei log provenienti da AVS
func AVS(ctx context.Context, logfile string) (err error) {

	start := time.Now()
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

	csvWriter := csv.NewWriter(gw)
	csvWriter.Comma= ';'
	// Scrive headers.
	gw.Write([]byte("#Log AVS prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte(avsheader + "\n"))
	


	// Apri file zippato in memoria
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file AVS %s, %s\n", logfile, err.Error())
		return
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	var s []string
	//var topic string
	n := 0
	for scan.Scan() {
		n++

		// AVS non ha header e quindi non lo salto
		line := scan.Text()

		// Verifica che logfile sia di tipo AVS.
		// if !isAVS.MatchString(line) {
		// 	err := fmt.Errorf("Error logfile %s non di tipo AVS", logfile)
		// 	return err
		// }

		// fmt.Println(line) // debug

		_, s, err = elaboraAVS2(ctx, line)
		if err != nil {
			log.Printf("Error Impossibile elaborare fruzione per record: %s", s)
		}

		// if len(s) < 2 {
		// 	continue
		// }

		// Scrive dati.
		err := csvWriter.Write(s)
		if err != nil {
			log.Printf("ERROR Impossibile srivere: %s\n", err.Error())
		}
		// justString := strings.Join(s, ";")
		// fmt.Println(justString)
		// gw.Write([]byte(justString + "\n"))

		//	fmt.Println(s[:]) // debug
		//records = append(records, strings.Join(s, ","))
	}

	// Scrive footer.
	gw.Write([]byte("#Numero di records: " + strconv.Itoa(n) + "\n"))
	gw.Close()

	// err = ioutil.WriteFile(newFile, []byte(justString), 0644)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

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

	fmt.Println(time.Since(start))
	return err
}

func elaboraAVS(ctx context.Context, line string) (topic string, result []string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Il separatore per i log AVS è "|"
	s := strings.Split(line, "|")

	t, err := time.ParseInLocation(timeAVSFormat, s[1], loc)
	if err != nil {
		log.Println(err.Error())
	}

	ora := t.UTC().Hour()
	minuto := t.UTC().Minute()

	// calcola a quale quartodora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15

	quartooraStr := strconv.Itoa(quartoora)

	//epoch := t.Format(time.RFC1123Z)
	//Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)

	// Crea il campo giornoq per integrare i log al quarto d'ora.
	giornoq := t.UTC().Format("20060102") + "q" + quartooraStr

	// idvideoteca := s[5]

	// ! OFFUSCAMENTO CAMPI SENSIBILI

	// Effettua hash della mail dell'utente.
	s[10], err = hasher.StringSumWithSalt(s[10], salt)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	// Effettua hash del cli utente.
	s[2], err = hasher.StringSumWithSalt(s[2], salt)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	e := strings.Join(s, `;`)

	result = append(result, giornoq, e)

	return giornoq, result, err
}

func elaboraAVS2(ctx context.Context, line string) (topic string, result []string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Elimina i divisori di campo nel campo delle mail secondarie
	// Se esistono virgolette nel record allora ci sono mail multiple
	if strings.Contains(line, `"`) {
		ll := strings.Split(line,`"`)
		if strings.Contains(ll[1],"|") {
			ll[1] = strings.Replace(ll[1], "|", " ",-1)
		}
		line = strings.Join(ll,"")
	}
	


	// Il separatore per i log AVS è "|"
	s := strings.Split(line, "|")

	t, err := time.ParseInLocation(timeAVSFormat, s[1], loc)
	if err != nil {
		log.Println(err.Error())
	}

	ora := t.UTC().Hour()
	minuto := t.UTC().Minute()

	// calcola a quale quartodora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15

	quartooraStr := strconv.Itoa(quartoora)

	//epoch := t.Format(time.RFC1123Z)
	//Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)

	// Crea il campo giornoq per integrare i log al quarto d'ora.
	giornoq := t.UTC().Format("20060102") + "q" + quartooraStr

	// idvideoteca := s[5]

	// ! OFFUSCAMENTO CAMPI SENSIBILI

	// Effettua hash della mail dell'utente.
	s[10], err = hasher.StringSumWithSalt(s[10], salt)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	// Gestione account secondari
	s[11], err = hasher.StringSumWithSalt(s[11], salt)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	// Effettua hash del cli utente.
	s[2], err = hasher.StringSumWithSalt(s[2], salt)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	//e := strings.Join(s, `;`)
	//Prepend field
	result = append([]string{giornoq}, s...)

	return giornoq, result, err
}
