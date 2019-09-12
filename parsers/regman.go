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
	"runtime"
	"strings"
	"sync"

	"github.com/axamon/hermes/idvideoteca"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

const headerregman = "cpeid;tgu;trap_timestamp;deviceid;devicetype;originipaddress;avgsskbps;bufferingduration;errordesc;errorreason;eventname;linespeedkbps;maxsschunkkbps;maxsskbps;minsskbps;videoduration;videoposition;videotype;videourl;eventtype;fwversion;networktype;ra_version;service_id_version"
const timeRegmanFormat = "2006-01-02 15:04:05"

// var isREGMAN = regexp.MustCompile(`(?m)^.*deviceid.*$`)

// NGASPLock gestisce l'accesso simultaneo alla scrittura sul file di output.
var NGASPLock sync.Mutex

var wgNGASP sync.WaitGroup

var writerchannel = make(chan *string, 1)

// REGMAN è il parser delle trap provenienti da REGMAN.
func REGMAN(ctx context.Context, logfile string, maxNumRoutines int) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Utilizzerà il massimo dei processori disponibili meno uno.
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	done := make(chan bool)

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
			case row := <-writerchannel:
				gw.Write([]byte(*row))
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
		switch {
		case numRoutines > maxNumRoutines:
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

	cpeid :=s[0]
	tgu :=s[1]
	trap_timestamp :=s[2]
	deviceid :=s[3]
	devicetype :=s[4]
	// mode :=s[5]
	originipaddress :=s[6]
	// averagebitrate :=s[7]
	avgsskbps :=s[8]
	bufferingduration :=s[9]
	// callerclass :=s[10]
	// callerrorcode :=s[11]
	// callerrormessage :=s[12]
	// callerrortype :=s[13]
	// callurl :=s[14]
	errordesc :=s[15]
	errorreason :=s[16]
	eventname :=s[17]
	// levelbitrates :=s[18]
	linespeedkbps :=s[19]
	maxsschunkkbps :=s[20]
	maxsskbps :=s[21]
	minsskbps :=s[22]
	// streamingtype :=s[23]
	videoduration :=s[24]
	videoposition :=s[25]
	// videotitle :=s[26]
	videotype :=s[27]
	videourl :=s[28]
	eventtype :=s[29]
	fwversion :=s[30]
	networktype :=s[31]
	ra_version :=s[32]
	// update_time :=s[33]
	// trap_provider :=s[34]
	// mid :=s[35]
	// service_id :=s[36]
	service_id_version :=s[37]
	// date_rif :=s[38]
	// video_provider :=s[39]
	// max_upstream_net_latency :=s[40]
	// min_upstream_net_latency :=s[41]
	// avg_upstream_net_latency :=s[42]
	// max_downstream_net_latency :=s[43]
	// min_downstream_net_latency :=s[44]
	// avg_downstream_net_latency :=s[45]
	// max_platform_latency :=s[46]
	// min_platform_latency :=s[47]
	// avg_platform_latency :=s[48]
	// packet_loss :=s[49]
	// preloaded_app_v :=s[50]
	


	// t, err := time.ParseInLocation(timeRegmanFormat, s[2], loc)
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// recupera ip cliente

	// ! OFFUSCAMENTO CAMPI SENSIBILI

	// Effettue hash ip pubblico cliente.
	originipaddressHASH, err := hasher.StringSumWithSalt(originipaddress, salt)

	// Effettue hash del cli cliente.
	tguHASH, err := hasher.StringSumWithSalt(tgu, salt)


	var resutlt []string
	//Prepend field
	result := append(resutlt,
		cpeid,
		tguHASH, // campo hashato
		trap_timestamp,
		deviceid,
		devicetype,
		originipaddressHASH, // campo hashato
		avgsskbps,
		bufferingduration,
		errordesc,
		errorreason,
		eventname,
		linespeedkbps,
		maxsschunkkbps,
		maxsskbps,
		minsskbps,
		videoduration,
		videoposition,
		videotype,
		videourl,
		eventtype,
		fwversion,
		networktype,
		ra_version,
		service_id_version,
		idv) // aggiunge campo con idunicovideoteca

	recordready := strings.Join(result, ";") + "\n"

	// Scrive dati.
	//err = csvWriter.Write(result)
	// NGASPLock.Lock()
	// gw.Write([]byte(recordready))
	// // gw.Flush()
	// NGASPLock.Unlock()
	writerchannel <- &recordready

	runtime.Gosched()
	return err
}
