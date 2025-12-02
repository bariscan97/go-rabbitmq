package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bariscan97/go-rabbitmq/event"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	conn := event.NewConnection("amqp://guest:guest@localhost:5672")
	defer conn.Close()

	consumer := event.NewConsumer(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	for _, s := range os.Args[1:] {
		log.Printf("Binding to severity: %s", s)
		
		go func(severity string) {
			err := consumer.Listen(ctx, "logs_direct", "direct", "", severity, func(body []byte) error {
				log.Printf(" [x] %s: %s", severity, body)
				return nil
			})
			if err != nil {
				log.Printf("Failed to listen for %s: %v", severity, err)
			}
		}(s)
	}

	<-ctx.Done()
}
