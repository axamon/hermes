package consumalog

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaLocalConsumer consuma i messaggi in kafka.
func KafkaLocalConsumer(ctx context.Context, topic string, oldoffset int64) (data []byte, offset int64, err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	partition := 0

	conn, _ := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()
	conn.Seek(oldoffset, 0)

	batch := conn.ReadBatch(10e3, 10e6) // fetch 10KB min, 1MB max
	defer batch.Close()

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b))
	}

	return b, batch.Offset(), err
}
