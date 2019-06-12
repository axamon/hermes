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
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/axamon/hermes/hasher"

	"github.com/axamon/hermes/zipfile"
)

var isCDN = regexp.MustCompile(`(?s)^\[.*\]\t[0-9]+\t\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\t[A-Z_]+\/\d{3}\t\d+\t[A-Z]+\t.*$`)

const seed = "vvkidtbcjujhgffbjnvrngvrinvufjkvljreucecvfcj"

// CDN è il parser dei log provenienti dalla Content Delivery Network
func CDN(logfile string) (err error) {

	fmt.Println(logfile)
	ctx := context.TODO()

	// Apri file zippato in memoria
	reader, err := zipfile.ReadAll(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		//return
	}

	// Riconosci tipo di file è veramente CDN
	scan := bufio.NewScanner(reader)
	for scan.Scan() {
		line := scan.Text()
		//fmt.Println(line)
		//if !isCDN.MatchString(line) {
		//	err := fmt.Errorf("Error logfile %s non di tipo CDN: %s", logfile, line)
		//	return err
		//}

		// Splitta ogni linea
		s := strings.Split(line, "\t")
		//Sfmt.Println(s)
		ip := s[2]
		ipHashed, err := hasher.StringSumWithSeed(ip, seed)
		if err != nil {
			log.Printf("Error Imposibile effettuare hashing %s\n", err.Error())
		}
		s[2] = ipHashed
		fmt.Println(s[:])
	}

	return err
}
