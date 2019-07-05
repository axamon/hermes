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

// KafkaLocalProducer3 produce messaggi in kafka.
func KafkaLocalProducer3(ctx context.Context, logfile string) (err error) {

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

	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", "logs", 0)
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
				// Partition: 0,
				// Topic:     topic,
				Value: []byte(line)},
		)
	}

	return

}

// KafkaLocalProducer2 produce messaggi in kafka.
func KafkaLocalProducer2(ctx context.Context, logfile string) (err error) {

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

	partition := 0

	n := 0
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		n++

		topic := strings.Split(line, ",")[0]

		conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)
		if err != nil {
			log.Printf("Error impossibile connettersi: %s\n", err.Error())
			return err
		}
		// fmt.Println(topic, line)

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		conn.WriteMessages(
			kafka.Message{
				// Partition: 0,
				// Topic:     topic,
				Value: []byte(line)},
		)
		conn.Close()

	}
	log.Printf("Prodotti %d records", n)
	return

}

var writers = make(map[string]*kafka.Writer)
var records = make(map[string][]string)
var canale = make(chan (*string))
var nlog int

// KafkaLocalProducer produce messaggi in kafka.
func KafkaLocalProducer(ctx context.Context, logfile string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Se non riesce a scrivere su Kafka procede senza andare in panico.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("Error impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		return err
	}

	r := bytes.NewReader(content)

	scan := bufio.NewScanner(r)

	nlog := 0
	// Produce record in kafka.
	go func() {
		for scan.Scan() {
			line := scan.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			canale <- &line
		}
	}()

	select {
	case <-canale:
		if len(canale) >= 100 {
			record := <-canale
			elabora(ctx, record)
		}

	default:
		record := <-canale
		elabora(ctx, record)
	}

	fmt.Println(nlog)
	return
}

func elabora(ctx context.Context, record *string) {
	nlog++
	topic := strings.Split(*record, ",")[0]

	if _, ok := writers[topic]; ok == false {
		writers[topic] = kafka.NewWriter(kafka.WriterConfig{Brokers: []string{"localhost:9092"}, Topic: topic})
		defer writers[topic].Close()
	}

	records[topic] = append(records[topic], *record)
	if len(records) >= 100 {
		for _, line := range records[topic] {
			err := writers[topic].WriteMessages(ctx, kafka.Message{Value: []byte(line)})
			if err != nil {
				log.Printf("Error Impossibile produrre record in kafka\n")
			}
		}
	}

	return
}
