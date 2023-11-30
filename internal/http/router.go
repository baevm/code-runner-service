package http

import (
	coderunnerHandler "code-runner-service/internal/handlers/coderunner"
	coderunnerSvc "code-runner-service/internal/services/coderunner"

	"github.com/labstack/echo/v4"
)

func (s *Server) LoadRoutes(e *echo.Echo) {
	crSvc := coderunnerSvc.New(s.mq, s.logger)

	ch := coderunnerHandler.New(s.logger, crSvc)

	e.POST("/run_code", ch.RunCodeHandler)
}
