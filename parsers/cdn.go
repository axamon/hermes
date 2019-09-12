// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers

import (
	"github.com/axamon/hermes/idvideoteca"
	"runtime"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

const cdnheader = "timestamp;hashfruizione;clientip;idvideoteca;status;tts[nanosecondi];bytes[bytes]"
const timeCDNFormat = "[02/Jan/2006:15:04:05.000+000]"

//var isCDN = regexp.MustCompile(`(?s)^\[.*\]\t[0-9]+\t\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\t[A-Z_]+\/\d{3}\t\d+\t[A-Z]+\t.*$`)

// CDNLock gestisce l'accesso simultaneo alla scrittura sul file di output.
var CDNLock sync.Mutex

var wgCDN sync.WaitGroup

var writerchannelcdn = make(chan string, 1)

// CDN è il parser dei log provenienti dalla Content Delivery Network
func CDN(ctx context.Context, logfile string, maxNumRoutines int) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Utilizzerà il massimo dei processori disponibili meno uno.
	runtime.GOMAXPROCS(runtime.NumCPU()-1)

	done := make(chan bool)

	// Apre nuovo file per salvare dati elaborati.
	newFile := strings.Split(logfile, ".csv.gz")[0] + ".offuscato.csv.gz"

	f, err := os.Create(newFile)
	if err != nil {
		log.Println(err.Error())
	}

	gw := gzip.NewWriter(f)
	defer gw.Close()

	// Scrive headers.
	//gw.Write([]byte("#Log CDN prodotto da piattaforma Hermes Copyright 2019 alberto.bregliano@telecomitalia.it\n"))
	gw.Write([]byte(cdnheader + "\n"))

	go func() {
		for {
			select {
			case row := <- writerchannelcdn:
					gw.Write([]byte(row))
			case <-done:
				return
			}
		}
	}()

	// Apre file zippato in memoria
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		return
	}

	r := bytes.NewReader(content)

	
	n := 0
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		n++
		line := scan.Text()

		numRoutines := runtime.NumGoroutine()
		wgCDN.Add(1)
		switch  {
		case  numRoutines > maxNumRoutines:
			ElaboraCDN(ctx, &line, gw)
		default:
			go ElaboraCDN(ctx, &line, gw)
		}
	}

	wgCDN.Wait()
	done <- true

	// Scrive footer.
	//gw.Write([]byte("#Numero di records: " + strconv.Itoa(n) + "\n"))
	gw.Flush()
	gw.Close()
	
	return err
}





// ElaboraCDN crea il file csv compresso con i campi sensibili offuscati.
func ElaboraCDN(ctx context.Context, line *string, gw *gzip.Writer) (err error) {
	
	ctx, cleanUP := context.WithCancel(ctx)
	defer cleanUP()

	defer wgCDN.Done()

	// Splitta la linea nei suoi fields,
	// il separatore per i log CDN è il tab: \t
	s := strings.Split(*line, "\t")


	timestamp := s[0]
	
	tts := s[1]

	bytes := s[4]

	// Recupera l'ip del cliente.
	clientip := s[2]

	// crea un idv vuoto
	var idv string
	idv, err = idvideoteca.Find(s[6])
	if err != nil {
		return err
	}

	

	// Converte il timestamp del log.
	// t, err := time.Parse(timeCDNFormat, s[0]) // UTC
	// if err != nil {
	// 	log.Println(err.Error())
	// }


	// Recupera lo status HTTP del chunk.
	status := s[3]

	// Recupera lo user agent del cliente.
	ua := s[8]

	

	// Crea l'idfruzione univoco del cliente.
	// cliemtip NON è ancora hashato.
	Hashfruizione, err := hasher.StringSum(clientip + idv + ua)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	

	// if len(str) < 2 {
	// 	return fmt.Errorf("Record troppo corto: %v", str)
	// }


	// ! OFFUSCAMENTO IP PUBBLICO CLIENTE
	// s[2] è l'ip pubblico del cliente da offuscare
	clientiphash, err := hasher.StringSumWithSalt(clientip, salt)
	if err != nil {
		log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
	}

	var str []string
	str = append(str, timestamp, Hashfruizione, clientiphash, idv, status, tts, bytes)

	record := strings.Join(str, ";") + "\n"
	//cdnrecords = append(cdnrecords, strings.Join(str, ";"))
	
	writerchannelcdn <- record

	runtime.Gosched()
	return err
}
