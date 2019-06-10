// il client serve a riconoscere la creazione di nuovi log
// tramite un watchdog
// a riconocere il tipo di log e a parsarlo di conseguenza
// Una volta sanitarizzato il log viene inviato a una destinazione
// prefissata

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var userid = flag.String("user", "pippo", "username")
var password = flag.String("pass", "pippo", "password")
var remoteAddr = flag.String("r", "http://127.0.0.1:8080", "default http://127.0.0.1:8080")
var file = flag.String("file", "", "no default")

type info struct {
	Name string `json: name`
	Data string `json: data`
}

func main() {

	cfx := context.Background()

	flag.Parse()

	var remoteURL string
	fmt.Println(*remoteAddr)
	remoteURL = *remoteAddr + "/upload"

	filedainviare := *file

	err := upload(remoteURL, filedainviare)
	if err != nil {
		log.Println(err.Error())
	}
}

func upload(url string, filedainviare string) (err error) {
	file, err := os.Open(filedainviare)
	if err != nil {
		log.Printf("impossible aprire file: %s errore: %s\n", filedainviare, err.Error())
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	if err != nil {
		log.Println(err.Error())
	}

	kvPairs, err := json.Marshal(info{Name: filedainviare, Data: encoded})

	//fmt.Printf("Sending JSON string '%s'\n", string(kvPairs))

	// Send request to OP's web server
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(kvPairs))
	if err != nil {
		log.Printf(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	//Aggiunge sicurezza
	req.SetBasicAuth(*userid, *password)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	//body, err := ioutil.ReadAll(resp.Body)

	//fmt.Println("Response: ", string(body))
	return
}
