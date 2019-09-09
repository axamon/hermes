// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers

import (
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// file del timestamp di roma embeddado in base64
const rome = "VFppZjIAAAAAAAAAAAAAAAAAAAAAAAAHAAAABwAAAAAAAACsAAAABwAAAA2AAAAAmzj4cJvVzOCcxcvwnbcAYJ6J/nCfoBzgoGCl8KF+rWCiXDdwo0waYMhsNfDM50sQzakXkM6CdODOokMQz5I0EM/jxuDQbl6Q0XIWENJM0vDTPjGQ1EnSENUd93DWKZfw1uuAkNgJlhD5M7Xw+dnE4Psc0nD7ubTw/Py0cP2ZlvD+5dDw/4KzcADFsvABYpVwApxacANCd3AEhXbwBSuT8AZuk3AHC3XwCEU68AjrV/AKLldwCss58AwOOXAMqxvwDeTg8A6K/fAPzf1wEHQacBGt33ASU/xwEs6X8BNNRBAUM/qQFSPrkBYT3JAXA82QF/O+kBjjr5AZ06CQGsORkBu8vRAcrK4QHZyfEB6MkBAffIEQIGxyECFcYxAiTFQQIzxFECQsNhAlHCcQJgwYECcFQ5An9TSQKOUlkCnVFpAqxQeQK7T4kCyk6ZAtlNqQLoTLkC90vJAwZK2QMV3ZEDJytBAzPbsQNFKWEDUdnRA2MngQNv1/EDgblJA43WEQOft2kDq9QxA721iQPKZfkD27OpA+hkGQP5sckEBmI5BBhDkQQkYFkENkGxBEJeeQRUP9EEYFyZBHI98QR+7mEEkDwRBJzsgQSuOjEEuuqhBMzL+QTY6MEE6soZBPbm4QUIyDkFFXipBSbGWQUzdskFRMR5BVF06QViwpkFb3MJBYFUYQWNcSkFn1KBBatvSQW9UKEFygERBdtOwQXn/zEF+UzhBgX9UQYX3qkGI/txBjXcyQZB+ZEGU9rpBmCLWQZx2QkGfol5Bo/XKQach5kGrdVJBrqFuQbMZxEG2IPZBuplMQb2gfkHCGNRBxUTwQcmYXEHMxHhB0RfkQdREAEHYvFZB28OIQeA73kHjQxBB57tmQerCmEHvOu5B8mcKQfa6dkH55pJB/jn+QAgECAQIBAgECAQIBAwQBAwQBAwECBAMEAwQDBAIEAwQDBAMEAwQDBAMEAwQDBAMEAwQDBAMEAwIFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgUGBQYFBgAAC7QAAAAAHCABBAAADhAACQAADhAACQAAHCABBAAAHCABBAAADhAACUxNVABDRVNUAENFVAAAAAABAQEBAAAAAAABAQpDRVQtMUNFU1QsTTMuNS4wLE0xMC41LjAvMwo="
var romebytes, _ = base64.StdEncoding.DecodeString(rome)
var loc, _ = time.LoadLocationFromTZData("Europe/Rome", romebytes)

// ! Necessita del file timelocation nel sistema operativo. 
// ! Su windows può essere un problema perchè non sempre è presente.
// var loc, _ = time.LoadLocation("Europe/Rome")

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

func giornoq(t time.Time) string {
	ora := t.UTC().Hour()
	minuto := t.UTC().Minute()

	// calcola a quale quartodora appartiene il dato.
	quartoora := ((ora * 60) + minuto) / 15
	quartooraStr := strconv.Itoa(quartoora)

	// Crea il campo giornoq per integrare i log al quarto d'ora.
	giornoq := t.UTC().Format("20060102") + "q" + quartooraStr

	return giornoq
}
