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

	e.GET("money-manager/v1/users/:user_id/balance", h.GetBalance)
	e.POST("money-manager/v1/users", h.AddFunds)
}
