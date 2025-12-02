package event

import (
	"context"
	"log"
)

type Consumer struct {
	conn *Connection
}

type Handler func(body []byte) error

func NewConsumer(conn *Connection) *Consumer {
	return &Consumer{
		conn: conn,
	}
}

func (c *Consumer) Listen(ctx context.Context, exchange, exchangeType, queueName, routingKey string, handler Handler) error {
	ch, err := c.conn.GetConnection().Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if exchange != "" {
		err = ch.ExchangeDeclare(
			exchange,
			exchangeType,
			true,  // durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return err
		}
	}

	q, err := ch.QueueDeclare(
		queueName,
		false, // durable 
		false, // delete 
		queueName == "",  // exclusive 
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	if exchange != "" {
		err = ch.QueueBind(
			q.Name,
			routingKey,
			exchange,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		false, // auto-ack 
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages in %s. To exit press CTRL+C", q.Name)
	<-ctx.Done()

	return nil
}