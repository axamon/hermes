package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/axamon/hermes/consumalog"
)

func main() {
	ctx := context.Background()

	topic := os.Args[1]
	offset, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Printf("ERROR impossibile trasformare in int: %s\n", err.Error())
	}

	offset64 := int64(offset)

	_, newoffset, err := consumalog.KafkaLocalConsumer(ctx, topic, offset64)

	if err != nil {
		log.Printf("ERROR impossibile consumare: %s\n", err.Error())
	}

	fmt.Println(newoffset)

}
