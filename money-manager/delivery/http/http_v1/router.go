package http_v1

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "money-manager/money-manager/delivery/docs"
	"money-manager/money-manager/usecase"
)

// @title         Money Manager API
// @name          Service for management of user's money
// @contact.email gusiatnikovanatalia@gmail.com
// @version       1.0
// @license.name  free-to-use-license

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

	e.GET("reports", h.MakeReport)

	e.GET("healthcheck", h.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
