package coderunner

import (
	"code-runner-service/internal/containers"
	"code-runner-service/internal/mqueue"

	"go.uber.org/zap"
)

type Service struct {
	log       *zap.SugaredLogger
	dockerCli *containers.Client
	producer  *mqueue.Producer
}

func New(mq *mqueue.MQueue, log *zap.SugaredLogger) *Service {
	client, _ := containers.New()

	producer, _ := mqueue.NewProducer(mq.Conn)

	go startCodeRunnerConsumer(log)

	return &Service{
		log:       log,
		dockerCli: client,
		producer:  producer,
	}
}
