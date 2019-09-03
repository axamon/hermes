package manifest

import
(
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var isDash = regexp.MustCompile(`(?m)\.mpd$`)
var isSS = regexp.MustCompile(`(?m)^.*\.ism`)


// Find restituisce la url del manifest per la url passata come argomento.
func Find(rawurl string) (urlmanifest string, err error) {

	decodedurl, err := url.QueryUnescape(rawurl)
	if err != nil {
		log.Println(err.Error())
	}

	rawurlbyte := []byte(decodedurl)

	// gestione dei dash
	if isDash.Match(rawurlbyte)  {

		urlmanifest = decodedurl

		_, err = http.Head(urlmanifest)
		if err != nil {
			log.Println(err.Error())
		}

	}

	// gestione degli SmoothStreaming
	if isSS.Match(rawurlbyte) {
	urlmanifest = isSS.FindString(decodedurl)+"/Manifest"

		_, err = http.Head(urlmanifest)
		if err != nil {
			log.Println(err.Error())
		}
	}


return urlmanifest, err
}