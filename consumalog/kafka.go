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

package consumalog

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaLocalConsumer consuma i messaggi in un kafka locale.
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

// KafkaRemoteConsumer consuma i messaggi in un kafka remoto.
func KafkaRemoteConsumer(ctx context.Context, remoteserver, topic, gruppoid string) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := kafka.NewReader(kafka.ReaderConfig{Topic: topic, GroupID: gruppoid, Brokers: []string{remoteserver}})

	defer r.CommitMessages(ctx)
	defer r.Close()
	defer log.Println("\n", r.Offset())

	for {
		messaggio, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println(err.Error())
			break
		}

		fmt.Println(string(messaggio.Value))
		r.CommitMessages(ctx)
	}
	// conn, err := kafka.DialLeader(ctx, "tcp", remoteserver, topic, partition)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	// defer conn.Close()
	// conn.Seek(oldoffset, 0)

	// defer fmt.Printf("Offset: %v\n", conn.Offset())

	// for {
	// 	messaggio, err := conn.ReadMessage(10e6)
	// 	if err != nil {
	// 		log.Panicln(err.Error())
	// 	}
	// 	fmt.Println(string(messaggio.Value))
	// }

	// batch := conn.ReadBatch(10e3, 10e6) // fetch 10KB min, 1MB max
	// defer batch.Close()

	// b := make([]byte, 10e3) // 10KB max per message

	// r := bytes.NewReader(b)

	// scan := bufio.NewScanner(r)

	// //var topic string
	// for scan.Scan() {

	// 	// AVS non ha header e quindi non lo salto
	// 	line := scan.Text()

	// 	// Scrive dati.
	// 	fmt.Println(line)

	// }

	return err
}
