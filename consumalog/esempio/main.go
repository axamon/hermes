package main

import (
	"context"
	"fmt"
	"log"

	"github.com/axamon/hermes/consumalog"
)

func main() {
	ctx := context.Background()
	_, offset, err := consumalog.KafkaLocalConsumer(ctx)

	if err != nil {
		log.Printf("ERROR impossibile consumare: %s\n", err.Error())
	}

	fmt.Println(offset)

}
