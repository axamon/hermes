package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/zipfile"
)

func main() {
	content, err := zipfile.ReadAll(os.Args[1])

	if err != nil {
		log.Printf("KO")
	}

	scan := bufio.NewScanner(content)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println(line)
	}
}
