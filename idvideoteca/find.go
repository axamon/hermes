package idvideoteca

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)(/)\d{8,8}(/)`)

// var onlyNum = regexp.MustCompile(`[^0-9]+`)

// Find trova l'id videoteca nella stringa passata come argomento se esiste
// altrimenti riporta "NON DISPONIBILE".
func Find(rawurl string) (idvideoteca string, err error) {

	if strings.Contains(rawurl, "%") {
		rawurl, err = url.QueryUnescape(rawurl)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}
	}

	// Cerca la prima corrispondenza
	idvideoteca = re.FindString(rawurl)

	// idvideoteca = onlyNum.ReplaceAllString(element, "")

	idvideoteca = strings.ReplaceAll(idvideoteca, "/", "")

	return idvideoteca, err
}
