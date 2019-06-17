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

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/axamon/hermes/hasher"
)

// seed per avere risutati hash personalizzati
const seed = "vvkidtbcjujhgffbjnvrngvrinvufjkvljreucecvfcj"

var timout = flag.Int64("timeout", 3, "tempo massimo per effettuare upload di un file")
var userid = flag.String("user", "pippo", "username")
var password = flag.String("pass", "pippo", "password")
var remoteAddr = flag.String("r", "127.0.0.1:8080", "default 127.0.0.1:8080")
var file = flag.String("file", "", "no default")

type info struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Hash string `json:"hash"`
}

func main() {

	flag.Parse()

	// Crea il contesto.
	ctx := context.Background()

	filedainviare := *file

	fi, err := os.Stat(filedainviare)
	if err != nil {
		log.Println(err.Error())
	}

	// in MB
	sizefile := fi.Size() / 1024 / 1024
	if sizefile > 100 {
		log.Printf("Error File %s troppo grande: %v\n", filedainviare, sizefile)
		return
	}

	// Verifica che il server remoto sia raggiungibile
	_, err = net.DialTimeout("tcp", *remoteAddr, time.Duration(3*time.Second))
	if err != nil {
		log.Printf("Server remoto non raggiungibile, error: %s\n", err.Error())
		return
	}

	// fmt.Println(*remoteAddr) // debug
	remoteURL := "http://" + *remoteAddr + "/upload" // ! TODO CAMBIARE IN HTTPS

	timeout := time.Duration(*timout) * time.Second

	err = upload(ctx, remoteURL, filedainviare, timeout)
	if err != nil {
		log.Println(err.Error())
	}
}

func upload(ctx context.Context, url, filedainviare string, timeout time.Duration) (err error) {

	// ! timeout Massimo tempo per terminare upload di un file
	ctx, cancelUpload := context.WithTimeout(ctx, timeout)
	defer cancelUpload()

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		return
	default:
		// Apre file da inviare.
		file, err := os.Open(filedainviare)
		if err != nil {
			log.Printf("impossible aprire file: %s errore: %s\n", filedainviare, err.Error())
		}
		defer file.Close()

		// Legge il file in memoria.
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err.Error())
		}

		// Encoda il file in base64.
		encoded := base64.StdEncoding.EncodeToString(data)
		if err != nil {
			log.Println(err.Error())
		}

		// Calcola hash del file con seed.
		hash, err := hasher.FileWithSeed(filedainviare, seed)
		if err != nil {
			log.Printf("ERROR Impossibile ricavare hash del file %s: %s\n", filedainviare, err.Error())
		}

		// log.Println(hash) // debug

		// Effettua il marshalling in json dai dati secondo il type info.
		kvPairs, err := json.Marshal(info{Name: filedainviare, Data: encoded, Hash: hash})

		// fmt.Printf("Sending JSON string '%s'\n", string(kvPairs)) // debug

		// Crea la web request in POST aggiungendo il file encodato.
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(kvPairs))
		if err != nil {
			log.Printf(err.Error())
		}

		// Aggiunge il contesto alla richiesta.
		req.WithContext(ctx)

		// Aggiunge header per processare json.
		req.Header.Set("Content-Type", "application/json")

		// Aggiunge sicurezza
		req.SetBasicAuth(*userid, *password)

		// Crea client http
		client := &http.Client{}

		// Invia la web request.
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error(), req.Context().Value("NomeFile"))
			return err
		}

		// Chiude il body della web response come da specifica.
		//defer resp.Body.Close()

		//body, err := ioutil.ReadAll(resp.Body)
		//fmt.Println("Response: ", string(body))

		switch resp.StatusCode == 200 {
		case true:
			log.Printf("INFO File %s trasferito correttamente\n", filedainviare)
		case false:
			log.Printf("ERROR Trasferimento del file %s non riuscito: %s\n", filedainviare, resp.Status)
		}

	}
	return err

}
