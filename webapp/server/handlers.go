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
		log.Fatal("Decoding error", err, err.Error())
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
	case true:
		// log.Printf("Bella prova zi! %s trasferito bene. Gli hash coincidono.\n", filename)
		log.Printf("INFO Salvato file %s, scritti: %d bytes", filename, n)
	}

	return

}
