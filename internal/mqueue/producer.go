package mqueue

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Ch *amqp091.Channel
}

func NewProducer(conn *amqp091.Connection) (*Producer, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return &Producer{
		Ch: ch,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return p.Ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
}
