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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/axamon/hermes/hasher"
)

func checkErr(msg string, err error) {
	if err != nil {
		log.Println(msg, err.Error())
	}
}

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Versione: 2.0\nAutore: Alberto Bregliano\n")
}

func upload(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// In caso di panico recupera senza killare il server.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}

	}()

	select {
	case <-ctx.Done():
		log.Printf("ERROR Timeout per invio raggiunto: %s\n", ctx.Err().Error())
		http.Error(w, ctx.Err().Error(), http.StatusRequestTimeout)
	default:

		// Recupera i dati user e pass di autorizzazione.
		auth := r.Header.Get("Authorization")

		// Se i dati di autorizzazione sono assenti chiude.
		if !strings.HasPrefix(auth, "Basic ") {
			log.Print("Invalid authorization:", auth)
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}

		// Recupera i dati di autorizzazione.
		up, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			log.Print("authorization decode error:", err)
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}

		// Verifica i dati di autorizzazione.
		if string(up) != userPass {
			log.Print("invalid username:password: ", string(up))
			http.Error(w, http.StatusText(unauth), unauth)
			return
		}
		//io.WriteString(w, "Goodbye, World!")
		//log.Println(r.Method)

		// Crea una istanza di info per salvare i dati in arrivo.
		element := info{}

		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR Impossibile leggere il corpo della richiesta: %s\n", err.Error())
		}

		// Salva dentro element i dati.
		err = json.Unmarshal(jsn, &element)
		if err != nil {
			log.Printf("ERROR Impossibile decodificare: %s\n", err.Error())
		}

		// fmt.Println(element) // debug

		// Assegna valori alle tre variabili recuperandole da element.
		filename, encoded, hashreceived := element.Name, element.Data, element.Hash

		// Decodifica i dati del file.
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		checkErr("ERROR Problema nel decoding: ", err)

		// Crea il file dove salvare i dati.
		f, err := os.Create("./" + filename)
		checkErr("ERROR Problema nel creare il file: ", err)
		defer f.Close()

		// Scrive nel file il contenuto decodificato.
		n, err := f.Write(decoded)
		checkErr("ERROR Impossibile scerivere nel file: ", err)

		// Forza chiusura del file per eseguire verifica di checksum.
		f.Close()

		// log.Println(hashreceived) // debug

		// Crea un hash di tutto il contenuto del file.
		hash, err := hasher.FileWithSeed(filename, seed)
		if err != nil {
			log.Printf("ERROR Impossibile ricavare hash del file %s: %s\n", filename, err.Error())
		}

		// log.Println(hash) // debug

		// Se l'hash creato equivale a quello ricevuto bene, altrimenti va in errore.
		switch hashreceived == hash {
		case false:
			log.Printf("ERROR trasferimento di %s non riuscito, hash non corrispondono.\n", filename)
			http.Error(w, http.StatusText(500), 500)
			w.Write([]byte("Errore nel trasferimento, hash non corrispondono"))
		case true:
			w.Write([]byte("Trasferimento OK, hash corrispondono"))
			log.Printf("INFO Salvato file %s con successo, scritti: %d bytes\n", filename, n)
		}
	}
	return

}
