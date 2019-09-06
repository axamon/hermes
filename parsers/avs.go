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

	// Sistema le linee con più account di posta.
	line = GestisciMailMultiple(line)

	// Il separatore per i log AVS è "|"
	s := strings.Split(line, "|")

	if len(s) != 18 {
		err = fmt.Errorf("Errore: Linea contiene un numero di campi diverso da 18: %s", s)
		log.Println(err)
		return err
	}

	// Crea IDVAS come hash di avs.TGU + avs.CPEID + avs.IDVIDEOTECA non modificati
	rawIDAVS := s[2] + s[3] + s[5]
	IDAVS, err := hasher.StringSum(rawIDAVS)

	// Considera i timestamp in orario locale non UTC
	t, err := time.ParseInLocation(timeAVSFormat, s[1], loc) // impostato in common.go
	if err != nil {
		log.Fatalf("Impossibile recuperare file per fuso orario locale: %s", err.Error())
	}

	// Calcolo il quarto d'ora di riferimento
	giornoq := giornoq(t)

	// ! OFFUSCAMENTO CAMPI SENSIBILI

	// Effettua hash della mail dell'utente.
	s[10], err = hasher.StringSumWithSalt(s[10], salt)

	// Gestione account secondari
	s[11], err = hasher.StringSumWithSalt(s[11], salt)

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
	}

	// Gestione CPEID s[3]
	// Ci sono CPEID che sono numeri enormi con la virgola
	// se li trovo ci faccio un hash
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
	AVSLock.Lock()
	gw.Write([]byte(recordready))
	AVSLock.Unlock()

	return err
}

// GestisciMailMultiple permette di effettuare hashing di email cliente
// multiple nei log AVS.
func GestisciMailMultiple(line string) string {
	if strings.Contains(line, `"`) {
		ll := strings.Split(line, `"`)
		if strings.Contains(ll[1], "|") {
			ll[1] = strings.Replace(ll[1], "|", " ", -1)
		}
		return strings.Join(ll, "")
	}
	return line
}

// Invia i records su kafka locale.
//err = inoltralog.LocalKafkaProducer(ctx, topic, records)
// err = inoltralog.RemoteKafkaProducer(ctx, "52.157.136.139:9092", topic, records)
// if err != nil {
// 	log.Printf("Error Impossibile salvare su kafka: %s\n", err.Error())
// }
