package coderunner

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CodeRunRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func (h *Handler) RunCodeHandler(c echo.Context) error {
	req := new(CodeRunRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := h.coderunnerSvc.RunCode(req.Language, req.Code)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
