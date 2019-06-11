// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package zipfile

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

// ReadAll legge il file zippato passato come parametro e restituisce
// un io.Reader e un eventuale errore.
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

	return gr, err
}
