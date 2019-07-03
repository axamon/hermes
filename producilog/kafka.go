package producilog

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/axamon/hermes/zipfile"
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

// KafkaLocalProducer produce messaggi in kafka.
func KafkaLocalProducer(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	//var topic string

	conn, err := kafka.DialContext(ctx, "tcp", "localhost:9092")
	if err != nil {
		log.Printf("Error impossibile connettersi: %s\n", err.Error())
		return err
	}
	defer conn.Close()

	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		topic := strings.Split(line, ",")[0]

		fmt.Println(topic, line)

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		conn.WriteMessages(
			kafka.Message{
				Partition: 0,
				Topic:     topic,
				Value:     []byte(line)},
		)
	}

	return

}
