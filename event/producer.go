package event

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn *Connection
	ch   *amqp.Channel
}

func NewProducer(conn *Connection) (*Producer, error) {
	p := &Producer{
		conn: conn,
	}
	if err := p.setupChannel(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Producer) setupChannel() error {
	if p.ch != nil && !p.ch.IsClosed() {
		return nil
	}

	ch, err := p.conn.GetConnection().Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		ch.Close()
		return fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	p.ch = ch
	return nil
}

func (p *Producer) Publish(exchange, routingKey string, body []byte) error {
	if err := p.setupChannel(); err != nil {
		return err
	}

	confirms := p.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	select {
	case confirmed := <-confirms:
		if confirmed.Ack {
			log.Printf("Message published to %s:%s", exchange, routingKey)
			return nil
		}
		return fmt.Errorf("failed to publish message: nack received")
	case <-ctx.Done():
		return fmt.Errorf("failed to publish message: timeout")
	}
}

func (p *Producer) PublishRPC(exchange, routingKey, replyTo, correlationId string, body []byte) error {
	if err := p.setupChannel(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          body,
			CorrelationId: correlationId,
			ReplyTo:       replyTo,
		},
	)
}

func (p *Producer) Close() error {
	if p.ch != nil {
		return p.ch.Close()
	}
	return nil
}
