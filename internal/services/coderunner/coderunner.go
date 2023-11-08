package coderunner

import "code-runner-service/internal/containers"

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) RunCode(lang string, code string) (string, error) {
	return containers.RunCodeContainer(lang, code)
}
