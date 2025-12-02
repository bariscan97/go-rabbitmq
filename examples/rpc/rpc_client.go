package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bariscan97/go-rabbitmq/event"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibRPC(n int) (res int, err error) {
	conn := event.NewConnection("amqp://guest:guest@localhost:5672")
	defer conn.Close()

	for i := 0; i < 10; i++ {
		if conn.GetConnection() != nil && !conn.GetConnection().IsClosed() {
			break
		}
		time.Sleep(1 * time.Second)
	}

	ch, err := conn.GetConnection().Channel()
	if err != nil {
		return 0, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return 0, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return 0, err
	}

	corrId := randomString(32)

	producer, err := event.NewProducer(conn)
	if err != nil {
		return 0, err
	}

	err = producer.PublishRPC(
		"",          // exchange
		"rpc_queue", // routing key
		q.Name,      // replyTo
		corrId,      // correlationId
		[]byte(strconv.Itoa(n)),
	)
	if err != nil {
		return 0, err
	}

	for d := range msgs {
		if d.CorrelationId == corrId {
			res, err = strconv.Atoi(string(d.Body))
			if err != nil {
				return 0, err
			}
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := 30
	if len(os.Args) > 1 {
		var err error
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Invalid number: %v", err)
		}
	}

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibRPC(n)
	if err != nil {
		log.Fatalf("Failed to handle RPC: %v", err)
	}

	log.Printf(" [.] Got %d", res)
}
