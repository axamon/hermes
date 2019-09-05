// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/axamon/hermes/parsers"
)

func parsa(v parsers.Videoteca) {
	fmt.Println(v.GetIDVideoteca())
	fmt.Println(v.GetManifestURL())
}

func main() {

	vod := parsers.Vod{RawURL: "http%3A%2F%2Fvodabr.cb.ticdn.it%2Fvideoteca2%2FV3%2FFilm%2F2017%2F06%2F50670127%2FSS%2F11473278%2F11473278_SD.ism%2FManifest%23https%3A%2F%2Flicense.cubovision.it%2FLicense%2Frightsmanager.asmx"}

	urlbrutta := parsers.RawURL{URL: "http%3A%2F%2Fvodabr.cb.ticdn.it%2Fvideoteca2%2FV3%2FFilm%2F2017%2F06%2F50670127%2FSS%2F11473278%2F11473278_SD.ism%2FManifest%23https%3A%2F%2Flicense.cubovision.it%2FLicense%2Frightsmanager.asmx"}

	
	parsa(vod)

	parsa(urlbrutta)
	
}





