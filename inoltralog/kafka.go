// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
