package main

import (
	"code-runner-service/config"
	"code-runner-service/internal/http"
	"code-runner-service/internal/mqueue"
	"code-runner-service/pkg/logger"
	"fmt"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	err = config.Load()

	if err != nil {
		log.Fatal(err)
	}

	// err := containers.PullImages()

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	mq, err := mqueue.Connect(
		config.Get().RabbitMQ.Username,
		config.Get().RabbitMQ.Password,
		fmt.Sprintf("%s:%s", config.Get().RabbitMQ.Host, config.Get().RabbitMQ.Port),
		config.Get().RabbitMQ.VHost)

	if err != nil {
		log.Fatal(err)
	}

	defer mq.Conn.Close()

	srv := http.New(log, mq)
	srv.Run(config.Get().Server.Host, config.Get().Server.Port)
}
