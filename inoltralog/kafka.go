// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package inoltralog

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// VerificaLocalKafka verifica se l'instanza è raggiungibile.
func VerificaLocalKafka(ctx context.Context) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Sceglie il topic su cui scirvere.
	topic := "logs"

	// Seclie la partizione Kafka su cui scrivere.
	partition := 0

	// Configura la connessione.
	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)
	defer conn.Close()
	if err != nil {
		log.Printf("Error impossibile aprire connessione a kafka\n")
	}

	return err
}

// VerificaRemoteKafka verifica se l'instanza è raggiungibile.
func VerificaRemoteKafka(ctx context.Context, remotekafkaserver string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Sceglie il topic su cui scirvere.
	topic := "logs"

	// Seclie la partizione Kafka su cui scrivere.
	partition := 0

	// Configura la connessione.

	conn, err := kafka.DialLeader(ctx, "tcp", remotekafkaserver, topic, partition)
	defer conn.Close()
	if err != nil {
		log.Printf("Error impossibile aprire connessione a kafka: %s\n", err.Error())
	}

	return err
}

// LocalKafkaProducer invia records a una istanza kafka locale.
func LocalKafkaProducer(ctx context.Context, topic string, s []string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w := kafka.NewWriter(kafka.WriterConfig{Brokers: []string{"localhost:9092"}, Topic: topic})
	defer w.Close()

	// Se non riesce a scrivere su Kafka procede senza andare in panico.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	nlog := 0
	// Produce record in kafka.
	for _, line := range s {
		nlog++
		err := w.WriteMessages(ctx, kafka.Message{Value: []byte(line)})
		if err != nil {
			log.Printf("Error Impossibile produrre record in kafka\n")
		}
	}

	log.Printf("Prodotti %d logs", nlog)

	return err
}

// RemoteKafkaProducer invia records a una istanza kafka locale.
func RemoteKafkaProducer(ctx context.Context, remotekafkaserver, topic string, s []string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w := kafka.NewWriter(kafka.WriterConfig{Brokers: []string{remotekafkaserver}, Topic: topic})
	defer w.Close()

	// Se non riesce a scrivere su Kafka procede senza andare in panico.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	nlog := 0
	// Produce record in kafka.
	for _, line := range s {
		nlog++
		err := w.WriteMessages(ctx, kafka.Message{Value: []byte(line)})
		if err != nil {
			log.Printf("Error Impossibile produrre record in kafka\n")
		}
	}

	log.Printf("Prodotti %d logs", nlog)

	return err
}
