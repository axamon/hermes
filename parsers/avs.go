// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

const avsheader = "giornoq;device;timestamp;tgu;cpeid;attivita;idvideoteca;standard;metodopagamento;ignoto;ignoto1;mailcliente;ignoto3;servizio;costruttore;tipoprog;case;rete;num;IDAVS"
const timeAVSFormat = "2006-01-02T15:04:05"

var isAVS = regexp.MustCompile(`(?m)^.*\|.*\|.*$`)

var avsrecords []string
var avsRecords sync.Mutex

var wgAVS sync.WaitGroup

// AVS è il parser dei log provenienti da AVS
func AVS(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	start := time.Now()

	// Apre nuovo file per salvare dati elaborati.
	newFile := strings.Split(logfile, ".csv.gz")[0] + ".offuscato.csv.gz"

	// fmt.Println(newFile)

	f, err := os.Create(newFile)
	if err != nil {
		return err
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	csvWriter := csv.NewWriter(gw)
	csvWriter.Comma = ';'
	// Scrive headers.
	//gw.Write([]byte("#Log AVS prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte(avsheader + "\n"))

	//var s []string
	//var topic string

	// Apri file zippato in memoria
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file AVS %s, %s\n", logfile, err.Error())
		return
	}

	r := bytes.NewReader(content)

	n := 0
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		n++

		wgAVS.Add(1)
		line := scan.Text()

		go ElaboraAVS(ctx, line, gw)
	}

	wgAVS.Wait()

	// Scrive footer.
	//gw.Write([]byte("#Numero di records: " + strconv.Itoa(n) + "\n"))
	// Garantisce che i dati vengano scritti tutti sul csv.
	gw.Flush()
	gw.Close()

	fmt.Println("Impiegato: ", time.Since(start))
	return err
}

// ElaboraAVS Crea il file csv compresso con i campi sensibili offuscati.
func ElaboraAVS(ctx context.Context, line string, gw *gzip.Writer) (err error) {

	defer wgAVS.Done()

	if strings.Contains(line, `"`) {
		ll := strings.Split(line, `"`)
		if strings.Contains(ll[1], "|") {
			ll[1] = strings.Replace(ll[1], "|", " ", -1)
		}
		line = strings.Join(ll, "")
	}

	// Il separatore per i log AVS è "|"
	s := strings.Split(line, "|")

	if len(s) != 18 {
		log.Fatal("Errore", s)
	}

	// Crea IDVAS come hash di avs.TGU + avs.CPEID + avs.IDVIDEOTECA non modificati
	rawIDAVS := s[2] + s[3] + s[5]
	IDAVS, err := hasher.StringSum(rawIDAVS)

	// Considera i timestamp in orario locale non UTC
	t, err := time.ParseInLocation(timeAVSFormat, s[1], loc) // impostato in common.go
	if err != nil {
		log.Println(err.Error())
	}

	// Calcolo il quarto d'ora di riferimento
	giornoq := giornoq(t)

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

	// Gestione TGU
	l := int(len(s[2]))
	switch {
	case l > 12:
		log.Printf("ERROR TGU maggiore di 12: %s", s)
	case l == 0: // se la TGU è vuota non effettua l'hashing
		break
	case l < 12:
		for l := len(s[2]); l <= 12; l++ {
			s[2] = "0" + s[2]
		}
		fallthrough
	case l == 12:
		s[2], err = hasher.StringSumWithSalt(s[2], salt)
		if err != nil {
			log.Printf("Error Hashing in errore: %s\n", err.Error())
		}
	}

	// Gestione CPEID s[3]
	if strings.Contains(s[3], ",") {
		s[3], err = hasher.StringSum(s[3])
	}

	//Prepend field
	result := append([]string{giornoq}, s...)

	// Aggiunge IDAVS alla fine
	result = append(result, IDAVS)

	recordready := strings.Join(result, ";") + "\n"

	// Scrive dati.
	//err = csvWriter.Write(result)
	avsRecords.Lock()
	gw.Write([]byte(recordready))
	avsRecords.Unlock()

	return err
}

// Invia i records su kafka locale.
//err = inoltralog.LocalKafkaProducer(ctx, topic, records)
// err = inoltralog.RemoteKafkaProducer(ctx, "52.157.136.139:9092", topic, records)
// if err != nil {
// 	log.Printf("Error Impossibile salvare su kafka: %s\n", err.Error())
// }
