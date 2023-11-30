package coderunner

import (
	"code-runner-service/config"
	"code-runner-service/internal/containers"
	"code-runner-service/internal/models"
	"code-runner-service/internal/mqueue"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const (
	EXCHANGE_NAME = "code-runner"
	QUEUE_NAME    = "run-code-requests"
	ROUTING_KEY   = "run-code-request"
)

func (s *Service) RunCode(ctx context.Context, lang string, code string) (string, error) {
	if _, isExist := containers.Images[lang]; !isExist {
		return "", errors.New("language not found")
	}

	codeReq := models.Code{
		Lang: lang,
		Body: code,
	}

	codeData, err := codeReq.ToJSON()

	if err != nil {
		return "", err
	}

	err = s.producer.Publish(ctx, EXCHANGE_NAME, ROUTING_KEY, codeData)

	if err != nil {
		return "", err
	}

	return "started container", nil
}

func startCodeRunnerConsumer(log *zap.SugaredLogger) {
	mq, err := mqueue.Connect(config.Get().RabbitMQ.Username,
		config.Get().RabbitMQ.Password,
		fmt.Sprintf("%s:%s", config.Get().RabbitMQ.Host, config.Get().RabbitMQ.Port),
		config.Get().RabbitMQ.VHost)

	if err != nil {
		log.Fatal(err)
	}

	defer mq.Conn.Close()

	consumer, err := mqueue.NewConsumer(mq.Conn, EXCHANGE_NAME, QUEUE_NAME, ROUTING_KEY)

	if err != nil {
		log.Fatalw("failed create consumer", "error", err)
	}

	defer consumer.Ch.Close()

	if err = consumer.ApplyQos(config.Get().Containers.MaxCount, 0, false); err != nil {
		log.Fatalw("failed to apply qos", "error", err)
	}

	msgs, err := consumer.Consume()

	if err != nil {
		log.Errorw("failed to consume message", "error", err)
	}

	forever := make(chan bool)

	for msg := range msgs {
		msg := msg

		go func() {
			time.Sleep(5 * time.Second)
			log.Infow("message received", "body", string(msg.Body))
			msg.Ack(false)
		}()
	}

	<-forever
}
