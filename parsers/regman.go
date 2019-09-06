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
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/axamon/hermes/idvideoteca"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

const headerregman = "giornoq;cpeid;tgu;trap_timestamp;deviceid;devicetype;mode;originipaddress;averagebitrate;avgsskbps;bufferingduration;callerclass;callerrorcode;callerrormessage;callerrortype;callurl;errordesc;errorreason;eventname;levelbitrates;linespeedkbps;maxsschunkkbps;maxsskbps;minsskbps;streamingtype;videoduration;videoposition;videotitle;videotype;videourl;eventtype;fwversion;networktype;ra_version;update_time;trap_provider;mid;service_id;service_id_version;date_rif;video_provider;max_upstream_net_latency;min_upstream_net_latency;avg_upstream_net_latency;max_downstream_net_latency;min_downstream_net_latency;avg_downstream_net_latency;max_platform_latency;min_platform_latency;avg_platform_latency;packet_loss;preloaded_app_v"

const timeRegmanFormat = "2006-01-02 15:04:05"

var isREGMAN = regexp.MustCompile(`(?m)^.*deviceid.*$`)

// NGASPLock gestisce l'accesso simultaneo alla scrittura sul file di output.
var NGASPLock sync.Mutex

var wgNGASP sync.WaitGroup

var writerchannel = make(chan string, 1)


// REGMAN è il parser delle trap provenienti da REGMAN.
func REGMAN(ctx context.Context, logfile string, maxNumRoutines int) (err error) {
	

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Utilizzerà il massimo dei processori disponibili meno uno.
	runtime.GOMAXPROCS(runtime.NumCPU()-1)

	done := make(chan bool)

	start := time.Now()

	// Apre nuovo file per salvare dati elaborati.
	newFile := strings.Split(logfile, ".csv.gz")[0] + ".offuscato.csv.gz"

	f, err := os.Create(newFile)
	if err != nil {
		return err
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	csvWriter := csv.NewWriter(gw)
	csvWriter.Comma = ';'

	// Scrive headers.
	//gw.Write([]byte("#Log REGMAN prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte(headerregman + "\n"))


	go func() {
		for {
			select {
			case row := <- writerchannel:
					gw.Write([]byte(row))
			case <-done:
				return
			}
		}
	}()

	// Apri file zippato in memoria.
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file NGASP %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	n := 0
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		n++

		// Salta header.
		if n == 1 {
			continue
		}

		line := scan.Text()

		numRoutines := runtime.NumGoroutine()
		wgNGASP.Add(1)
		switch  {
		case  numRoutines > maxNumRoutines:
			ElaboraREGMAN(ctx, &line, gw)
		default:
			go ElaboraREGMAN(ctx, &line, gw)
	}

	
	}

	wgNGASP.Wait()
	done <- true
	// defer close(writerchannel)

	// Scrive footer.
	//gw.Write([]byte("#Numero di records: " + strconv.Itoa(n) + "\n"))
	gw.Flush()
	gw.Close()

	fmt.Println("Impiegato: ", time.Since(start))
	return err
}

// ElaboraREGMAN crea il file csv compresso con i campi sensibili offuscati.
func ElaboraREGMAN(ctx context.Context, line *string, gw *gzip.Writer) (err error) {

	ctx, cleanUP := context.WithCancel(ctx)
	defer cleanUP()

	defer wgNGASP.Done()



	// ricerca le fruzioni nell'intervallo temporale richiesto
	// l'intervallo temporale inzia con l'inzio di una fruizione

	// Splitta la linea nei suoi campi.
	// Il separatore per i log REGMAN è ";"
	s := strings.Split(*line, ";")

	// crea un idv vuoto
	var idv string

	// Se è un VOD Estrae id videoteca univoco del vod
	if strings.Contains(strings.ToLower(s[27]), "vod") {
		idv, _ = idvideoteca.Find(s[28])
	}

	// Crea IDNGASP come hash di ngasp.TGU + ngasp.CPEID + ngasp.IDVIDEOTECA non modificati
	rawIDNGASP := s[0] + s[1] + idv
	IDNGASP, err := hasher.StringSum(rawIDNGASP)

	t, err := time.ParseInLocation(timeRegmanFormat, s[2], loc)
	if err != nil {
		log.Println(err.Error())
	}

	giornoq := giornoq(t)

	// recupera ip cliente

	// ! OFFUSCAMENTO CAMPI SENSIBILI

	// Effettue hash ip pubblico cliente.
	s[6], err = hasher.StringSumWithSalt(s[6], salt)
	if err != nil {
		log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
	}

	// Effettue hash del cli cliente.
	s[1], err = hasher.StringSumWithSalt(s[1], salt)
	if err != nil {
		log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
	}

	// Eliminazione campo titolo
	s[26] = "" // questo è il campo con il nome del film viene sostituito con idvideoteca
	s[36] = "" // nei games ci sono titoli che hanno apici
	s[33] = "" // a volte questo campo ha apici
	for n, l := range s {
		if strings.Contains(l, `'`) {
			fmt.Printf("Il record contiente caratteri non accettati: %d, %s\n", n, s)
		}
	}

	//Prepend field
	result := append([]string{giornoq}, s...)

	// Aggiunge IDNGASP alla fine
	result = append(result, IDNGASP)

	recordready := strings.Join(result, ";") + "\n"

	// Scrive dati.
	//err = csvWriter.Write(result)
	// NGASPLock.Lock()
	// gw.Write([]byte(recordready))
	// // gw.Flush()
	// NGASPLock.Unlock()
	writerchannel <- recordready

	runtime.Gosched()
	return err
}
