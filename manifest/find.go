// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manifest

import (
//	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var isDash = regexp.MustCompile(`(?m)\.mpd$`)
var isSS = regexp.MustCompile(`(?m)^.*\.ism`)

// Find restituisce la url del manifest per la url passata come argomento.
func Find(rawurl string) (urlmanifest string, err error) {

	// serve a rendere gestibili le url encodate
	if strings.Contains(rawurl, "%") {
		rawurl, err = url.QueryUnescape(rawurl)
		if err != nil {
			return "", err
		}
	}

	//rawurlbyte := []byte(decodedurl)

	// gestione dei dash
	if isDash.MatchString(rawurl) {
		urlmanifest = rawurl
	}

	// gestione degli SmoothStreaming
	if isSS.MatchString(rawurl) {
		urlmanifest = isSS.FindString(rawurl) + "/Manifest"
	}

	return urlmanifest, err
}

// IsManifestReachable esegue un curl HEAD sul manifest passato come argomento.
// func IsManifestReachable(urlmanifest string) error {
// 	_, err := http.Head(urlmanifest)
// 	return err
// }
