package main

import (
	"fmt"

	"github.com/axamon/hermes/parsers"
)

func main() {

	var v = new(parsers.Vod)

	v.RawURL = "http%3A%2F%2Fvodabr.cb.ticdn.it%2Fvideoteca2%2FV3%2FFilm%2F2017%2F06%2F50670127%2FSS%2F11473278%2F11473278_SD.ism%2FManifest%23https%3A%2F%2Flicense.cubovision.it%2FLicense%2Frightsmanager.asmx"

	id := v.GetIDVideoteca()
	m := v.GetManifestURL()

	v.IDVideoteca = id
	v.ManifestURL = m

	fmt.Println(id, m, v)
}
