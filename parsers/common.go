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
	"log"
	"net/url"
	"strings"
	"regexp"
	"os"
	"time"
)

var loc, _ = time.LoadLocation("Europe/Rome")

var salt = os.Getenv("salt")

// Utenti clienti TimVision.
type Utenti []Fruizioni

// Fruizioni archivia tuttue le fruizioni effettuate.
type Fruizioni struct {
	Hashfruizione map[string]bool
	Clientip      map[string]string
	Idvideoteca   map[string]string
	Idaps         map[string]string
	Edgeip        map[string]string
	Giorno        map[string]string
	Orario        map[string]string
	Details       map[string][]float64 `json:"-"`
}


// Per idvideoteca
var idV = regexp.MustCompile(`(?m)(/)\d{7,8}(/)`)

// Per manifest
var isDash = regexp.MustCompile(`(?m)\.mpd$`)
var isSS = regexp.MustCompile(`(?m)^.*\.ism`)

// Videoteca è una inerfaccia per i metodi di estrazione dati utili.
type Videoteca interface {
	GetIDVideoteca() string
	GetManifestURL() string
}

func extractIDVideoteca(rawurl string) string {
	var idvideoteca string

	if strings.Contains(rawurl, "%") {
		rawurlDecoded, err := url.QueryUnescape(rawurl)
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		rawurl = rawurlDecoded
	}

	// Cerca la prima corrispondenza
	idvideoteca = idV.FindString(rawurl)

	// idvideoteca = onlyNum.ReplaceAllString(element, "")

	idvideoteca = strings.Replace(idvideoteca, "/", "", 2)
	
	return idvideoteca
}


func extractManifest(rawurl string) (urlmanifest string, err error) {

		// serve a rendere gestibili le url encodate
		if strings.Contains(rawurl, "%") {
			rawurl, err = url.QueryUnescape(rawurl)
			if err != nil {
				log.Println(err.Error())
				return "", err
			}
		}
	
	
		// gestione dei dash
		if isDash.MatchString(rawurl) {
			urlmanifest = rawurl
		}
	
		// gestione degli SmoothStreaming
		if isSS.MatchString(rawurl) {
			urlmanifest = isSS.FindString(rawurl) + "/Manifest"
		}

	return urlmanifest, nil
}