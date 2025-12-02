package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	consumer := event.NewConsumer(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received signal, shutting down...")
		cancel()
	}()

	err := consumer.Listen(ctx, "logs", "fanout", "", "", func(body []byte) error {
		log.Printf(" [x] %s", body)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
