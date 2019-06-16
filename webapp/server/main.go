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
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

//const userPass = "pippo:pippo"
const unauth = http.StatusUnauthorized
const timedout = http.StatusRequestTimeout

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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	// userPass contiene le credenziali
	userPass = *userid + ":" + *password

	s := &http.Server{
		Addr: *port,
		//	Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/version", version)

	http.HandleFunc("/upload", upload)

	log.Panic(s.ListenAndServe())

}
