package mqueue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Message queue using RabbitMQ
type MQueue struct {
	Conn *amqp.Connection
}

type MQClient struct {
	Ch *amqp.Channel
}

func Connect(username, password, host, vhost string) (*MQueue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))

	if err != nil {
		return nil, err
	}

	return &MQueue{
		Conn: conn,
	}, nil
}

func (mq MQClient) Close() error {
	return mq.Ch.Close()
}
