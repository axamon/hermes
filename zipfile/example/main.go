package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/axamon/hermes/zipfile"
)

func main() {
	content, err := zipfile.ReadAll("20190607_03_00_vodabr.cb.log.gz")

	if err != nil {
		log.Printf("KO")
	}

	scan := bufio.NewScanner(content)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println(line)
	}
}
