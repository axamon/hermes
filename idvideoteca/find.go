package idvideoteca

import (
	"strings"
    "regexp"
    "fmt"
)

var re = regexp.MustCompile(`(?m)(/|%2F)\d{8,8}(/|%2F)`)

// Find trova l'id videoteca nella stringa passata come argomento se esiste
// altrimenti riporta "NON DISPONIBILE"
func Find(s string) (idvideoteca string, err error) {
    
	elements := re.FindAllString(s, 1)
	if len(elements) <1 {
	return "", fmt.Errorf("idvideoteca NON DISPONIBILE")
	}

	if strings.Contains(elements[0], "/") {

		idvideoteca = strings.Replace(elements[0], "/","",-1)
	}

	if strings.Contains(elements[0], "%2F") {
		idvideoteca = strings.Replace(elements[0], "%2F","",-1)
	}

	return idvideoteca, nil
	
}