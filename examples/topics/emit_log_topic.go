package main

import (
	"log"
	"os"
	"strings"

	"github.com/bariscan97/go-rabbitmq/event"
)

func main() {
	conn := event.NewConnection("amqp://guest:guest@localhost:5672")
	defer conn.Close()

	producer, err := event.NewProducer(conn)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	body := bodyFrom(os.Args)
	severity := severityFrom(os.Args)

	err = producer.Publish("logs_topic", severity, []byte(body))
	if err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}

	log.Printf(" [x] Sent %s: %s", severity, body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
