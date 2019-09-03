package main

import (
	"log"
	"os"
	"fmt"
	"github.com/axamon/hermes/manifest"
)

func main() {
	
manifesturl, err := manifest.Find(os.Args[1])
if err != nil {
	log.Println(err.Error())
}

fmt.Println(manifesturl)
}