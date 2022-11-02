package http_v1

import (
	"github.com/labstack/echo/v4"

	"money-manager/money-manager/usecase"
)

type ServerHandler struct {
	uc usecase.MoneyManagerUseCase
}

func NewServerHandler(e *echo.Echo, useCase usecase.MoneyManagerUseCase) {
	h := &ServerHandler{
		uc: useCase,
	}

	e.GET("/accounts/{account_id}/balance", h.GetFunds)
}

func (e *ServerHandler) GetFunds(eCtx echo.Context) error {
	return nil
}
