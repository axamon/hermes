package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	endpoint = "https://gegup.telecomitalia.local:8443/DiagnosticTool/api.php?method=DiagnosticTool&sincrono=N&format=json&tgu="
)

func main() {

	ctx := context.Background()

	tgu := os.Args[1]

	client := &http.Client{}
	url := endpoint + tgu

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("ERROR Impossibile creare richiesta: %s\n", err.Error())
	}

	req.WithContext(ctx)

	username := os.Getenv("DiagnosticToolUsername")
	password := os.Getenv("DiagnosticToolPassoword")

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR Impossibile inviare richiesta http: %s\n", err.Error())
	}
	defer resp.Body.Close()

	responsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR Impossibile leggere body reqest: %s\n", err.Error())
	}

	fmt.Println(string(responsBody))

}
