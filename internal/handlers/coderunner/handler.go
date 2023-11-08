package coderunner

import (
	"code-runner-service/internal/services/coderunner"

	"go.uber.org/zap"
)

type Handler struct {
	log           *zap.SugaredLogger
	coderunnerSvc *coderunner.Service
}

func New(log *zap.SugaredLogger, coderunnerSvc *coderunner.Service) *Handler {
	return &Handler{
		log:           log,
		coderunnerSvc: coderunnerSvc,
	}
}
