package main

import (
	"log"
	"strconv"

	"github.com/bariscan97/go-rabbitmq/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func main() {
	conn := event.NewConnection("amqp://guest:guest@localhost:5672")
	defer conn.Close()
	
	ch, err := conn.GetConnection().Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	producer, err := event.NewProducer(conn)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			if err != nil {
				log.Printf("Failed to convert body to integer: %v", err)
				d.Ack(false)
				continue
			}

			log.Printf(" [.] fib(%d)", n)
			response := fib(n)

			err = producer.PublishRPC(
				"",        // exchange
				d.ReplyTo, // routing key
				"",        // replyTo 
				d.CorrelationId,
				[]byte(strconv.Itoa(response)),
			)
			if err != nil {
				log.Printf("Failed to publish reply: %v", err)
			}

			d.Ack(false)
		}
	}()

	log.Printf(" [x] Awaiting RPC requests")
	<-forever
}
