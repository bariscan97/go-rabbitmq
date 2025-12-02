package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/bariscan97/go-rabbitmq/event"
)

func main() {
	conn := event.NewConnection("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	for i := 0; i < 10; i++ {
		if conn.GetConnection() != nil && !conn.GetConnection().IsClosed() {
			break
		}
		log.Println("Waiting for connection...")
		time.Sleep(1 * time.Second)
	}

	producer, err := event.NewProducer(conn)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	body := bodyFrom(os.Args)
	
	err = producer.Publish("logs", "", []byte(body))
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
