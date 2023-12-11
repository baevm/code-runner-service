package coderunner

import (
	"code-runner-service/config"
	"code-runner-service/internal/containers"
	"code-runner-service/internal/models"
	"code-runner-service/internal/mqueue"
	"code-runner-service/pkg/random"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

const (
	EXCHANGE_NAME = "code-runner"
	QUEUE_NAME    = "run-code-requests"
	ROUTING_KEY   = "run-code-request"
)

func (s *Service) AddCodeToQueue(ctx context.Context, lang string, code string) (string, error) {
	if _, isExist := containers.Images[lang]; !isExist {
		return "", errors.New("language not found")
	}

	requestId, err := random.String(12)

	if err != nil {
		return "", err
	}

	codeReq := models.Code{
		Lang:      lang,
		Body:      code,
		RequestId: requestId,
	}

	codeData, err := codeReq.ToJSON()

	if err != nil {
		return "", err
	}

	// add code task to queue
	err = s.producer.Publish(ctx,
		EXCHANGE_NAME,
		ROUTING_KEY,
		codeData,
		codeReq.RequestId,
	)

	if err != nil {
		return "", err
	}

	return codeReq.RequestId, nil
}

func startCodeRunnerConsumer(log *zap.SugaredLogger, dockerCli *containers.Client) {
	mq, err := mqueue.Connect(config.Get().RabbitMQ.Username,
		config.Get().RabbitMQ.Password,
		fmt.Sprintf("%s:%s", config.Get().RabbitMQ.Host, config.Get().RabbitMQ.Port),
		config.Get().RabbitMQ.VHost)

	if err != nil {
		log.Fatalw("failed to connect to rabbitmq", "error", err)
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

		// start new goroutine for each message
		// to process multiple messages at the same time
		go func() {
			defer msg.Ack(false)
			log.Infow("message received", "id", msg.MessageId)

			var code models.Code

			err := json.Unmarshal(msg.Body, &code)

			if err != nil {
				log.Errorw("failed to parse message", "error", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), config.Get().Containers.MaxTime)
			defer cancel()

			res, err := dockerCli.RunCodeContainer(ctx, code)

			if err != nil {
				log.Errorw("failed to run code container", "error", err)
				return
			}

			log.Infow("code container finished", "result", res)
		}()
	}

	<-forever
}
