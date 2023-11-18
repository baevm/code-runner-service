package coderunner

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type CodeRunRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type CodeRunResponse struct {
	Result string `json:"result"`
}

func (h *Handler) RunCodeHandler(c echo.Context) error {
	req := new(CodeRunRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.coderunnerSvc.RunCode(ctx, req.Language, req.Code)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CodeRunResponse{Result: res})
}
