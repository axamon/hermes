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
	"fmt"
	"io"
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
	fmt.Fprintf(w, "Versione: 2.0\nAutore: Alberto Bregliano")
}

func upload(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}

	}()

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
		log.Print("errore reading body ", err, err.Error())
		http.Error(w, http.StatusText(timedout), timedout)
		return
	}

	err = json.Unmarshal(jsn, &element)
	if err != nil {
		log.Fatal("Decoding error ", err, err.Error())
	}

	// fmt.Println(element)

	filename, encoded, hashreceived := element.Name, element.Data, element.Hash

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	checkErr("Problema nel decoding: ", err)

	//fmt.Println(filename, decoded)
	f, err := os.Create("./" + filename)
	checkErr("Problema nel creare il file: ", err)
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
		w.Write([]byte("Errore nel trasferimento di: %s, hash non corrispondono"))
	case true:
		// log.Printf("Bella prova zi! %s trasferito bene. Gli hash coincidono.\n", filename)
		log.Printf("INFO Salvato file %s, scritti: %d bytes", filename, n)
	}

	return

}
