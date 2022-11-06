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

	e.POST("money-manager/users/balance", h.AddFunds)
	e.GET("money-manager/users/balance", h.GetBalance)
	e.PATCH("money-manager/users/balance", h.DebitFunds)

	e.POST("money-manager/transfers", h.TransferFunds)

	e.POST("money-manager/reserves", h.AddReserve)

	//todo: path adn verb for reserve revoke
	e.POST("money-manager/reserves/revoke", h.RevokeReserve)
	e.PATCH("money-manager/reserves", h.AcceptReserve)
}
