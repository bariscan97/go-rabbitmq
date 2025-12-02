package main

import (
	"bytes"
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

	handler := func(body []byte) error {
		log.Printf("Received a message: %s", body)
		dotCount := bytes.Count(body, []byte("."))
		t := time.Duration(dotCount)
		time.Sleep(t * time.Second)
		log.Printf("Done")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	err := consumer.Listen(ctx, "", "", "task_queue", "", handler)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
}
