package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/zipfile"
)

func main() {
	ctx := context.Background()

	content, err := zipfile.ReadAll(ctx, os.Args[1])

	if err != nil {
		log.Printf("KO")
	}

	scan := bufio.NewScanner(content)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println(line)
	}
}
