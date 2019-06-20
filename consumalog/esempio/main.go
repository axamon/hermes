package main

import (
	"context"

	"github.com/axamon/hermes/consumalog"
)

func main() {
	ctx := context.Background()
	consumalog.KafkaLocalConsumer(ctx)
}
