package parsers

import (
	"log"
	"net/url"
	"regexp"
)

// Per idvideoteca
var idv = regexp.MustCompile(`(?m)(/|%2[fF])\d{8,8}(/|%2[fF])`)
var onlyNum = regexp.MustCompile(`[^0-9]+`)

// Per manifest
var isDash = regexp.MustCompile(`(?m)\.mpd$`)
var isSS = regexp.MustCompile(`(?m)^.*\.ism`)

type videoteca interface {
	GetIDVideoteca() string
	GetManifestURL() string
}

type Vod struct {
	RawURL      string
	IDVideoteca string
	ManifestURL string
}

func (v Vod) GetIDVideoteca() string {

	decodedurl, err := url.QueryUnescape(v.RawURL)
	if err != nil {
		log.Println(err.Error())
	}

	// Trova la prima corrispondenza a sinistra
	element := idv.FindString(decodedurl)

	idvideoteca := onlyNum.ReplaceAllString(element, "")

	return idvideoteca
}

func (v Vod) GetManifestURL() string {

	var urlmanifest string

	// serve a rendere gestibili le url encodate
	decodedurl, err := url.QueryUnescape(v.RawURL)
	if err != nil {
		log.Println(err.Error())
	}

	rawurlbyte := []byte(decodedurl)

	// gestione dei dash
	if isDash.Match(rawurlbyte) {
		urlmanifest = decodedurl
	}

	// gestione degli SmoothStreaming
	if isSS.Match(rawurlbyte) {
		urlmanifest = isSS.FindString(decodedurl) + "/Manifest"
	}

	return urlmanifest
}
