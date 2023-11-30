package mqueue

import (
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Ch           *amqp091.Channel
	queueName    string
	exchangeName string
	key          string
}

func NewConsumer(conn *amqp091.Connection, exchangeName, queueName, key string) (*Consumer, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(exchangeName, "direct", false, false, false, false, nil); err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)

	if err != nil {
		return nil, err
	}

	if err := ch.QueueBind(queueName, key, exchangeName, false, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		Ch:           ch,
		exchangeName: exchangeName,
		queueName:    queueName,
		key:          key,
	}, nil
}

func (c *Consumer) Consume() (<-chan amqp091.Delivery, error) {
	return c.Ch.Consume(c.queueName, c.exchangeName, false, false, false, false, nil)
}

func (c *Consumer) ApplyQos(count int, size int, global bool) error {
	return c.Ch.Qos(count, size, global)
}
