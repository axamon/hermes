// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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