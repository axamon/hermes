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
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/axamon/hermes/hasher"
)

// ElaboraCDN trasforma ogni record dei log.
func ElaboraCDN(ctx context.Context, line string) (s []string, err error) {

	// ricerca le fruzioni nell'intervallo temporale richiesto
	// l'intervallo temporale inzia con l'inzio di una fruizione

	// fmt.Println(line)
	if !isCDN.MatchString(line) {
		err := fmt.Errorf("Error record %s non di tipo CDN: %s", line)
		return nil, err
	}

	// Splitta la linea nei supi fields.
	// Il separatore per i log CDN è il tab: \t
	s = strings.Split(line, "\t")

	//parsa la URL nelle sue componenti
	u, err := url.Parse(s[6])
	if err != nil {
		log.Printf("Error nel parsing URL di: %s\n", line)
	}
	/* Urlschema := u.Scheme
	if Urlschema != "https" { //fa passare solo le URL richieste via WEB
		continue
	} */

	//converte i timestamp come piacciono a me
	t, err := time.Parse("[02/Jan/2006:15:04:05.000+000]", s[0])
	if err != nil {
		log.Println(err.Error())
	}

	ora := t.Hour()
	minuto := t.Minute()

	// calcola a quale quartodora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15

	quartooraStr := strconv.Itoa(quartoora)

	IDipq, _ := hasher.StringSum(s[2] + quartooraStr)

	//epoch := t.Format(time.RFC1123Z)

	Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)
	var speed, tts, bytes float64

	tts, err = strconv.ParseFloat(s[1], 8)
	if err != nil {
		log.Fatal(err.Error())
	}

	bytes, err = strconv.ParseFloat(s[4], 8)
	if err != nil {
		log.Fatal(err.Error())
	}

	speed = (bytes / tts)
	speedStr := fmt.Sprintf("%f", speed)
	clientip := s[2]
	status := s[3] //da usare per errori 40x e 50x
	ua := s[8]

	//fmt.Println(Urlschema)
	//Urlhost := u.Host
	Urlpath := u.Path
	//fmt.Println(Urlpath)
	//Urlquery := u.RawQuery
	//Urlfragment := u.Fragment
	pezziurl := strings.Split(Urlpath, "/")
	//fmt.Println(pezziurl)

	var idvideoteca, idaps, Hash string

	if ok := !strings.Contains(Urlpath, "video="); ok == true { //solo i chunk video

		return nil, nil
	}
	idvideoteca = pezziurl[6]
	//tipocodifica := pezziurl[7]
	idaps = pezziurl[8]
	//fmt.Println(idvideoteca)
	//encoding := pezziurl[10]
	//fmt.Println(encoding)
	//re := regexp.MustCompile(`QualityLevels\(([0-9]+)\)$`)
	//bitratestr := re.FindStringSubmatch(encoding)[1]
	//bitrate, _ := strconv.ParseFloat(bitratestr, 8)
	/* if err != nil {
		log.Fatal(err.Error())
	} */
	//bitrateMB := bitrate * bitstoMB
	Hash, err = hasher.StringSum(clientip + idvideoteca + ua)
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	//ingestafruizioni(Hash, clientip, idvideoteca, idaps, edgeip, giorno, orario, speed)

	s = append(s, Time, Hash, idaps, idvideoteca, status, speedStr, quartooraStr, IDipq)

	return s, err
}

func ElaboraREGMAN(ctx context.Context, line string) (s []string, err error) {

	// ricerca le fruzioni nell'intervallo temporale richiesto
	// l'intervallo temporale inzia con l'inzio di una fruizione

	// fmt.Println(line)
	/* if !isREGMAN.MatchString(line) {
		err := fmt.Errorf("Error record %s non di tipo REGMAN: %s", line)
		return nil, err
	} */

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

	IDipq, _ := hasher.StringSum(s[6] + quartooraStr)

	//epoch := t.Format(time.RFC1123Z)

	Time := t.Format("200601021504") //idem con patate questo è lo stracazzuto ISO8601 meglio c'è solo epoch
	//fmt.Println(Time)

	Hash, err := hasher.StringSum(s[0])
	if err != nil {
		log.Printf("Error Hashing in errore: %s\n", err.Error())
	}

	//ingestafruizioni(Hash, clientip, idvideoteca, idaps, edgeip, giorno, orario, speed)

	s = append(s, Time, Hash, IDipq, quartooraStr)

	return s, err
}
