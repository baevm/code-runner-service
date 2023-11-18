package main

import (
	"code-runner-service/config"
	"code-runner-service/internal/http"
	"code-runner-service/pkg/logger"
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

	srv := http.New(log)
	srv.Run(config.Get().Server.Host, config.Get().Server.Port)
}
