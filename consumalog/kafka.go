package consumalog

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaLocalConsumer consuma i messaggi in kafka.
func KafkaLocalConsumer(ctx context.Context) (data []byte, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	topic := "logs"
	partition := 0

	conn, _ := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	defer batch.Close()

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b))
	}

	return
}
