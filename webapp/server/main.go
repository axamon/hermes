// il server provvede un endpoint a cui il client può
// inviare i log sanitarizzati.

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
)

//const userPass = "pippo:pippo"
const unauth = http.StatusUnauthorized

var client *http.Client
var remoteURL string


var userid = flag.String("user", "pippo", "username")
var password = flag.String("pass", "pippo", "password")
var port = flag.String("port", ":8080", "default :8080")

type info struct {
	Name string `json: name`
	Data string `json: data`
}

var userPass string

func main() {
	flag.Parse()

	userPass = *userid+":"+*password
	
	s := &http.Server{
		Addr: *port,
		//	Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.HandleFunc("/saluto", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Ciao, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "Ciao, straniero")

	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Ciao, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(w, "Versione: 1.0\nAutore: Alberto Bregliano")

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

		filename := element.Name

		encoded := element.Data
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

		log.Printf("Salvato file %s, scritti: %d bytes", filename, n)

		return

	})

	log.Fatal(s.ListenAndServe())
}
