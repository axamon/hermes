package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/hermes/titolovod"
)

func main() {

	ctx := context.Background()

	//var result = new(titolovod.Response)
	var err error

	idvideoteca := os.Args[1]
	fmt.Printf("cerco: %s", idvideoteca)

	result, err := titolovod.Get(ctx, idvideoteca)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(result[0].ResultObj.Containers[0].Metadata.Title)
}
