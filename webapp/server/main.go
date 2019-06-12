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
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/axamon/hermes/hasher"
)

//const userPass = "pippo:pippo"
const unauth = http.StatusUnauthorized

var client *http.Client
var remoteURL string

var userid = flag.String("user", "pippo", "username")
var password = flag.String("pass", "pippo", "password")
var port = flag.String("port", ":8080", "default :8080")

type info struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Hash string `json:"hash"`
}

var userPass string

// seed per avere risutati hash personalizzati
const seed = "vvkidtbcjujhgffbjnvrngvrinvufjkvljreucecvfcj"

func main() {
	flag.Parse()

	// userPass contiene le credenziali
	userPass = *userid + ":" + *password

	s := &http.Server{
		Addr: *port,
		//	Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Ciao, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "Versione: 2.0\nAutore: Alberto Bregliano")

	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth, "Basic ") {
			log.Print("Invalid authorization:", auth)
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
		up, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			log.Print("authorization decode error:", err)
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
		if string(up) != userPass {
			log.Print("invalid username:password: ", string(up))
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
		io.WriteString(w, "Goodbye, World!")
		//log.Println(r.Method)

		element := info{}

		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("errore reading body", err, err.Error())
		}

		err = json.Unmarshal(jsn, &element)
		if err != nil {
			log.Fatal("Deconing error", err, err.Error())
		}

		// fmt.Println(element)

		filename := element.Name

		encoded := element.Data

		hashreceived := element.Hash

		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			log.Println(err.Error())
		}

		//fmt.Println(filename, decoded)
		f, err := os.Create("./" + filename)
		if err != nil {
			log.Println(err.Error())
		}
		defer f.Close()

		n, err := f.Write(decoded)
		if err != nil {
			log.Println(err.Error())
		}

		f.Close()

		// log.Println(hashreceived) // debug

		hash, err := hasher.FileWithSeed(filename, seed)
		if err != nil {
			log.Printf("ERROR Impossibile ricavare hash del file %s: %s\n", filename, err.Error())
		}

		// log.Println(hash) // debug

		switch hashreceived == hash {
		case false:
			log.Printf("Errore nel trasferimento di: %s, hash non corrispondono.\n", filename)
		case true:
			// log.Printf("Bella prova zi! %s trasferito bene. Gli hash coincidono.\n", filename)
			log.Printf("INFO Salvato file %s, scritti: %d bytes", filename, n)
		}

		return

	})

	log.Fatal(s.ListenAndServe())
}
