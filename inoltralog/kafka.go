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
	"time"

	"github.com/segmentio/kafka-go"
)

// LocalKafkaProducer invia records a una istanza kafka locale.
func LocalKafkaProducer(ctx context.Context, s []string) (err error) {

	// Se non riesce a scrivere su Kafka procede senza andare in panico.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

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

	// Imposta timeout per la scrittura sull'istanza kafka.
	err = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.Printf("Error Timeout connessione a kafka\n")
	}

	// Produce record in kafka.
	for _, line := range s {
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(line)},
		)
		if err != nil {
			log.Printf("Error produrre record in kafka\n")
		}
	}

	return err
}
