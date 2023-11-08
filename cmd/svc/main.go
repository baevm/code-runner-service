package main

import (
	"code-runner-service/config"
	"code-runner-service/internal/http"
	"code-runner-service/pkg/logger"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	err = config.Load()

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(config.Get().Containers.MaxTime)

	// err := containers.PullImages()

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	srv := http.New(logger)
	srv.Run(config.Get().Server.Host, config.Get().Server.Port)
}
