package event

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn     *amqp.Connection
	url      string
	err      chan error
	closed   bool
}

func NewConnection(url string) *Connection {
	c := &Connection{
		url: url,
		err: make(chan error),
	}
	go c.connect()
	return c
}

func (c *Connection) connect() {
	for {
		log.Printf("Connecting to RabbitMQ at %s...", c.url)
		conn, err := amqp.Dial(c.url)
		if err == nil {
			c.conn = conn
			c.closed = false
			log.Println("Connected to RabbitMQ")
			
			closeErr := make(chan *amqp.Error)
			c.conn.NotifyClose(closeErr)
			
			err := <-closeErr
			if err != nil {
				log.Printf("Connection closed: %v", err)
				c.err <- err
			} else {
				log.Println("Connection closed gracefully")
				return
			}
		} else {
			log.Printf("Failed to connect: %v", err)
		}

		if c.closed {
			return
		}

		log.Println("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func (c *Connection) Close() error {
	c.closed = true
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) GetConnection() *amqp.Connection {
	return c.conn
}
