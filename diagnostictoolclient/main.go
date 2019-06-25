package main

import (
        "fmt"
        "regexp"
        "crypto/tls"
        "io/ioutil"
        "log"
        "net/http"
        "os"
)

const (
        endpointTgu = "https://10.38.34.138:8443/DiagnosticTool/api.php?method=DiagnosticTool&sincrono=N&format=json&tgu="
        endpointEsito = "https://10.38.34.138:8443/DiagnosticTool/api.php?method=DiagnosticTool&sincrono=Y&format=json&cod_esito="
)


var re = regexp.MustCompile(`(?m)^\d+$`)

func main() {

        var endpoint string

        parametro := os.Args[1]

        if !re.MatchString(parametro) {
        endpoint = endpointEsito
        }


        if re.MatchString(parametro) {
        endpoint = endpointTgu
        }

        // Costringe il client ad accettare anche certificati https non validi
        // o scaduti.
        transCfg := &http.Transport{
                // Ignora certificati SSL scaduti.
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        client := &http.Client{Transport: transCfg}
        url := endpoint + parametro

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Printf("ERROR Impossibile creare richiesta: %s\n", err.Error())
        }


        username := os.Getenv("DiagnosticToolUsername")
        password := os.Getenv("DiagnosticToolPassoword")

        req.SetBasicAuth(username, password)

        resp, err := client.Do(req)
        if err != nil {
                log.Printf("ERROR Impossibile inviare richiesta http: %s\n", err.Error())
        }
        //defer resp.Body.Close()

        responsBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Printf("ERROR Impossibile leggere body reqest: %s\n", err.Error())
        }

        fmt.Println(string(responsBody))

}
