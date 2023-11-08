package http

import (
	"code-runner-service/internal/handlers/coderunner"

	"github.com/labstack/echo/v4"
)

func (s *Server) LoadRoutes(e *echo.Echo) {
	ch := coderunner.New()

	e.POST("/code", ch.RunCodeHandler)
}
