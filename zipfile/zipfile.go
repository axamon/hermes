package zipfile

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

// ReadAll legge il file zippato passato come parametro e restituisce
// il suo contenuto ed eventuale errore.
func ReadAll(zipFile string) (content io.Reader, err error) {

	// Apre il file zippato in lettura.
	f, err := os.Open(zipFile)
	defer f.Close()
	if err != nil {
		log.Printf("Errore Impossibile aprire il file %s: %s\n", zipFile, err.Error())
	}

	gr, err := gzip.NewReader(f)
	defer gr.Close()
	if err != nil {
		log.Printf("Errore Impossibile leggere il contenuto del file %s: %s\n", zipFile, err.Error())
	}

	// content, err = ioutil.ReadAll(gr)
	// if err != nil {
	//	log.Printf("Errore Impossibile prelevare contenuti del file file %s: %s\n", zipFile, err.Error())
	// }

	return gr, err
}
