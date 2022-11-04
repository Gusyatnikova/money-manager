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

	e.GET("money-manager/users/balance", h.GetBalance)
	e.POST("money-manager/users/balance", h.AddFunds)

	e.POST("money-manager/users/balance/reserve", h.AddReserve)
	e.DELETE("money-manager/users/balance/reserve", h.AcceptReserve)
}
