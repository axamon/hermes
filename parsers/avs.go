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

	"github.com/axamon/hermes/hasher"
	"github.com/axamon/hermes/zipfile"
)

var isAVS = regexp.MustCompile(`(?m)^.*\|.*\|.*$`)

// AVS è il parser dei log provenienti da AVS
func AVS(logfile string) (err error) {

	ctx := context.TODO()

	// Apri file zippato in memoria
	content, err := zipfile.ReadAllZIP(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file AVS %s, %s\n", logfile, err.Error())
		return
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	n := 0
	for scan.Scan() {
		n++
		line := scan.Text()

		// Verifica che logfile sia di tipo CDN.
		if !isAVS.MatchString(line) {
			err := fmt.Errorf("Error logfile %s non di tipo AVS", logfile)
			return err
		}

		fmt.Println(line) // debug

		s, err := ElaboraAVS(ctx, line)
		if err != nil {
			log.Printf("Error Impossibile elaborare fruzione per record: %s", s)
		}

		if len(s) < 2 {
			continue
		}

		// ! ANONIMIZZAZIONE CAMPI SENSIBILI

		mailutente := s[10]
		mailHashed, err := hasher.StringSumWithSeed(mailutente, seed)
		if err != nil {
			log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
		}
		s[10] = mailHashed

		utente := s[11]
		utenteHashed, err := hasher.StringSumWithSeed(utente, seed)
		if err != nil {
			log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
		}
		s[11] = utenteHashed

		fmt.Println(s[:])
	}

	fmt.Println(n)
	return err
}
