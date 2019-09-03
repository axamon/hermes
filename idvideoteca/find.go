package idvideoteca

import (
	"net/url"
	"log"
//	"strings"
    "regexp"
    "fmt"
)

var re = regexp.MustCompile(`(?m)(/|%2[fF])\d{8,8}(/|%2[fF])`)
var onlyNum = regexp.MustCompile(`[^0-9]+`)

// Find trova l'id videoteca nella stringa passata come argomento se esiste
// altrimenti riporta "NON DISPONIBILE".
func Find(rawurl string) (idvideoteca string, err error) {
	
	decodedurl, err := url.QueryUnescape(rawurl)
	if err != nil {
		log.Println(err.Error())
	}

	// Cerca la prima corrispondenza
	elements := re.FindAllString(decodedurl, 1)
	
	// Se non ci sono corrispondenze esce con errore
	if len(elements) <1 {
		return "", fmt.Errorf("idvideoteca NON DISPONIBILE")
	}


	idvideoteca = onlyNum.ReplaceAllString(elements[0], "")

	return idvideoteca, nil
	
}