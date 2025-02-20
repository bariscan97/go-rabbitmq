package main

import (
	"os"
    "github.com/bariscan97/go-rabbitmq/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	consumer, err := event.NewConsumer(connection)
	if err != nil {
		panic(err)
	}
	consumer.Listen(os.Args[1:])
}