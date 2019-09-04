package idvideoteca

import (
	"net/url"
	"regexp"
	"strings"
)

var idv = regexp.MustCompile(`(?m)(/)\d{7,8}(/)`)
// var isURLEncoded = regexp.MustCompile(`%`)

// var onlyNum = regexp.MustCompile(`[^0-9]+`)

// Find trova l'id videoteca nella stringa passata come argomento se esiste
// altrimenti riporta "NON DISPONIBILE".
func Find(rawurl string) (idvideoteca string, err error) {

	if strings.Contains(rawurl, "%") {
	// if isURLEncoded.MatchString(rawurl) {
		rawurl, err = url.QueryUnescape(rawurl)
		if err != nil {
			// log.Println(err.Error())
			return "", err
		}
	}

	// Cerca la prima corrispondenza
	idvideoteca = idv.FindString(rawurl)

	// idvideoteca = onlyNum.ReplaceAllString(element, "")

	idvideoteca = strings.Replace(idvideoteca, "/", "", 2)

	return idvideoteca, err
}
