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

	e.POST("users/balance", h.AddMoney)
	e.GET("users/balance", h.GetBalance)
	e.PATCH("users/balance", h.DebitMoney)

	e.PATCH("transfers", h.TransferMoney)

	e.POST("reserves/add", h.AddReserve)
	e.DELETE("reserves/cancel", h.CancelReserve)
	e.DELETE("reserves/accept", h.AcceptReserve)

	e.GET("healthcheck", h.HealthCheck)
}
