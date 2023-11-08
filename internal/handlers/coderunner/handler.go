package coderunner

import (
	"code-runner-service/internal/containers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

type CodeRunRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func (h *Handler) RunCodeHandler(c echo.Context) error {
	req := new(CodeRunRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := containers.RunCodeContainer(req.Language, req.Code)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
