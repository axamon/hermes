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
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/inoltralog"
	"github.com/axamon/hermes/zipfile"
)

var isREGMAN = regexp.MustCompile(`(?m)^.*\;.*\;.*\;.*\;.*\;.*$`)

// REGMAN è il parser delle trap provenienti da REGMAN.
func REGMAN(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// fmt.Println(logfile) // debug

	// Apri file zippato in memoria.
	content, err := zipfile.ReadAllZIP(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file REGMAN %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	var records []string
	n := 0
	for scan.Scan() {
		n++
		line := scan.Text()

		// Verifica che logfile sia di tipo regman.
		if !isREGMAN.MatchString(line) {
			err := fmt.Errorf("Error logfile %s non di tipo REGMAN", logfile)
			return err
		}

		s, err := elaboraREGMAN(ctx, line)
		if err != nil {
			log.Printf("Error Impossibile elaborare REGMAN record: %s", s)
		}

		if len(s) < 2 {
			continue
		}

		//! OFFUSCAMENTO CAMPI SENSIBILI
		// s[2] contiente ip pubblico cliente.
		s[2], err = hasher.StringSumWithSalt(s[2], salt)
		if err != nil {
			log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
		}

		// s[3] contiene il cli del cliente.
		s[3], err = hasher.StringSumWithSalt(s[3], salt)
		if err != nil {
			log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
		}
		//	fmt.Println(s[:]) // debug

		// Viene aggiunto a records un record con i campi individuati
		// separati da tab.
		records = append(records, strings.Join(s, "\t"))
	}

	// Scrive uno per uno su standard output i record offuscati.
	for _, record := range records {
		fmt.Println(record)
	}

	// Invia i records su kafka locale.
	err = inoltralog.LocalKafkaProducer(ctx, records)
	if err != nil {
		log.Printf("Error Impossibile salvare su kafka: %s\n", err.Error())
	}

	fmt.Println(n)
	return err
}

func elaboraREGMAN(ctx context.Context, line string) (s []string, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// ricerca le fruzioni nell'intervallo temporale richiesto
	// l'intervallo temporale inzia con l'inzio di una fruizione

	// fmt.Println(line)
	if !isREGMAN.MatchString(line) {
		err := fmt.Errorf("Error recordnon di tipo REGMAN: %s", line)
		return nil, err
	}

	// Splitta la linea nei supi fields.
	// Il separatore per i log REGMAN è ";"
	s = strings.Split(line, ";")

	t, err := time.Parse("2006-01-02 15:04:05", s[2])
	if err != nil {
		log.Println(err.Error())
	}

	ora := t.Hour()
	minuto := t.Minute()

	// calcola a quale quartodora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15

	quartooraStr := strconv.Itoa(quartoora)

	//IDipq, _ := hasher.StringSum(s[6] + quartooraStr)

	//epoch := t.Format(time.RFC1123Z)

	// Crea il campo giornoq per integrare i log al quarto d'ora.
	giornoq := []string{t.Format("20060102") + "q" + quartooraStr}

	//Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)

	// recupera ip cliente
	ipregman := s[6]
	cli := s[1]

	//ingestafruizioni(Hash, clientip, idvideoteca, idaps, edgeip, giorno, orario, speed)

	result := append(giornoq, "REG", ipregman, cli, s[0], s[5])

	return result, err
}
