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
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

const avsheader = "device;timestamp;tgu;cpeid;idvideoteca;costruttore;rete;num"
const timeAVSFormat = "2006-01-02T15:04:05"

var isAVS = regexp.MustCompile(`(?m)^.*\|.*\|.*$`)

// AVSLock gestisce l'accesso simultaneo alla scrittura sul file di output.
var AVSLock sync.Mutex

var wgAVS sync.WaitGroup

// AVS è il parser dei log provenienti da AVS
func AVS(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Apre nuovo file per salvare dati elaborati.
	newFile := strings.Split(logfile, ".csv.gz")[0] + ".offuscato.csv.gz"

	f, err := os.Create(newFile)

	gw := gzip.NewWriter(f)
	defer gw.Close()

	csvWriter := csv.NewWriter(gw)
	csvWriter.Comma = ';'

	// Scrive headers.
	//gw.Write([]byte("#Log AVS prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte(avsheader + "\n"))

	// Apre file zippato in memoria
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file AVS %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	n := 0
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		n++

		line := scan.Text()

		wgAVS.Add(1)
		go ElaboraAVS(ctx, line, gw)
	}

	wgAVS.Wait()

	// Scrive footer.
	//gw.Write([]byte("#Numero di records: " + strconv.Itoa(n) + "\n"))
	// Garantisce che i dati vengano scritti tutti sul csv.
	gw.Flush()
	gw.Close()

	return err
}

// ElaboraAVS Crea il file csv compresso con i campi sensibili offuscati.
func ElaboraAVS(ctx context.Context, line string, gw *gzip.Writer) (err error) {

	defer wgAVS.Done()


	// Il separatore per i log AVS è "|"
	s := strings.Split(line, "|")

	// Considera i timestamp in orario locale non UTC
	// t, err := time.ParseInLocation(timeAVSFormat, s[1], loc) // impostato in common.go
	// if err != nil {
	// 	log.Fatalf("Impossibile recuperare file per fuso orario locale: %s", err.Error())
	// }



	device := s[0]
	timestamp := s[1]
	tgu := s[2]
	cpeid := s[3]
	idvideoteca := s[5]
	costruttore := s[13]
	rete := s[16]
	num := s[17]

		// ! OFFUSCAMENTO CAMPI SENSIBILI

	

	// Gestione tgu s[2]
	tgu, err = hasher.StringSumWithSalt("tgu", salt)

	// Gestione CPEID s[3]
	// Ci sono CPEID che sono numeri enormi con la virgola
	// se li trovo ci faccio un hash
	if strings.Contains(cpeid, ",") {
		cpeid, err = hasher.StringSum(cpeid)
	}

	var result []string

	result = append(result, device, timestamp, tgu, cpeid, idvideoteca, costruttore, rete, num)
	recordready := strings.Join(result, ";") + "\n"

	// Scrive dati.
	//err = csvWriter.Write(result)
	AVSLock.Lock()
		gw.Write([]byte(recordready))
	AVSLock.Unlock()

	return err
}