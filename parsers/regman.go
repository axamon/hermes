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

var isREGMAN = regexp.MustCompile(`(?m)^\w+\;\d{12}\;[0-9-\s:\.]+\;\w+\;\w+\;\w+\;[0-9\.]{8,16}\;.*$`)

// REGMAN è il parser dei log provenienti di regman.
func REGMAN(logfile string) (err error) {

	fmt.Println(logfile)
	ctx := context.TODO()

	// Apri file zippato in memoria
	//reader, err := zipfile.ReadAll(ctx, logfile)

	content, err := zipfile.ReadAllZIP(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file REGMAN %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	n := 0
	for scan.Scan() {
		n++
		line := scan.Text()

		// Verifica che logfile sia di tipo regman.
		/* if !isREGMAN.MatchString(line) {
			err := fmt.Errorf("Error logfile %s non di tipo REGMAN", logfile)
			return err
		} */

		s, err := ElaboraREGMAN(ctx, line)
		if err != nil {
			log.Printf("Error Impossibile elaborare fruzione per record: %s", s)
		}

		if len(s) < 2 {
			continue
		}

		// ! ANONIMIZZAZIONE IP PUBBLICO CLIENTE
		ip := s[6]
		ipHashed, err := hasher.StringSumWithSeed(ip, seed)
		if err != nil {
			log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
		}
		s[6] = ipHashed

		fmt.Println(s[:])
	}

	fmt.Println(n)
	return err
}
