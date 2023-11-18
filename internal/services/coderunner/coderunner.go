package coderunner

import (
	"code-runner-service/internal/containers"
	"code-runner-service/internal/models"
	"context"
)

type Service struct {
	dockerCli *containers.Client
}

func New() *Service {
	client, _ := containers.New()

	return &Service{
		dockerCli: client,
	}
}

func (s *Service) RunCode(ctx context.Context, lang string, code string) (string, error) {
	codeReq := models.Code{
		Lang: lang,
		Body: code,
	}

	return s.dockerCli.RunCodeContainer(ctx, codeReq)
}
