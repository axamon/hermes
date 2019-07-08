package producilog

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
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

var writers = make(map[string]*kafka.Writer)
var records = make(map[string][]string)
var canale = make(chan *string, 100)
var nlog int
var wg sync.WaitGroup

// KafkaLocalProducer produce messaggi in kafka.
func KafkaLocalProducer(ctx context.Context, logfile string) (err error) {

	// Crea il contesto e la funzione di cancellazione.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	start := time.Now()

	// Se non riesce a scrivere su Kafka procede senza andare in panico.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	done := make(chan bool, 1)

	// Produce records in kafka.
	log.Println("Avvio select")
	// go func() {
	// 	defer log.Println("Select  finito")
	// 	for {
	// 		select {
	// 		case <-done:
	// 			return
	// 		default:
	// 			wg.Add(1)
	// 			elabora(ctx)
	// 		}
	// 	}
	// }()

	// Apre il file zippato e salva il contenuto in memoria.
	content, err := zipfile.ReadAllGZ(ctx, logfile)
	if err != nil {
		log.Printf("ERROR Impossibile leggere file CDN %s, %s\n", logfile, err.Error())
		return err
	}

	// Trasforma il contenuto in *Reader
	r := bytes.NewReader(content)

	// Crea uno scanner
	scan := bufio.NewScanner(r)

	fmt.Println("Ciclo Scan iniziato")
	startScan := time.Now()
	for scan.Scan() {
		// line := new(string) // si usa new per creare line nella heap
		line := scan.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}
		//	fmt.Println(*line)
		canale <- &line
		wg.Add(1)
		go elabora(ctx)
		//fmt.Println("linea caricata su canale")
	}
	fmt.Println("Ciclo Scan finito", time.Since(startScan))

	// Chiude il canale.
	close(canale)

	// Comunica che i cicli scan sono finiti.
	done <- true

	// Attende che tutte le elaborazioni siano finite.
	wg.Wait()

	// Mostra quanti records sono stati processati.
	fmt.Println(nlog)

	// Mostra quanto tempo è stato richiesto per terminare elaborazione.
	fmt.Println(time.Since(start))

	return
}

func elabora(ctx context.Context) {
	fmt.Println("Inizio Goroutine")
	defer wg.Done()

	for record := range canale {
		// fmt.Println("Elaboro linea")
		var topic string
		topic = strings.Split(*record, ",")[0]

		if _, ok := writers[topic]; ok == false {

			// time.Sleep(2 * time.Microsecond)
			writers[topic] = kafka.NewWriter(kafka.WriterConfig{Brokers: []string{"localhost:9092"}, Topic: topic})
			defer writers[topic].Close()
		}

		records[topic] = append(records[topic], *record)
		fmt.Println(len(records[topic]))
		_, isOpen := <-canale
		if len(records[topic]) >= 100 || isOpen == false {
			for _, line := range records[topic] {

				strings.Split(line, ",")
				// time.Sleep(2 * time.Microsecond)

				err := writers[topic].WriteMessages(ctx, kafka.Message{Value: []byte(line)})
				if err != nil {
					log.Printf("Error Impossibile produrre record in kafka\n")
				}
			}
			nlog++
		}
		fmt.Println(nlog)
	}

	runtime.Gosched()
	return
}
