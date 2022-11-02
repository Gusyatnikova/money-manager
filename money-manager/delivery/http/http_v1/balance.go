package http_v1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"money-manager/money-manager/entity"
)

type balance struct {
	CurAmount   uint64 `json:"current_amount"`
	AvailAmount uint64 `json:"available_amount"`
}

type userBalance struct {
	UserId string  `json:"user_id"`
	Ub     balance `json:"balance"`
}

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	usrParam := eCtx.Param("user_id")
	if !e.isValidUserId(usrParam) {
		return e.noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.GetBalance(): invalid userId"))
	}

	usr := entity.User{
		UserId: usrParam,
	}

	bal, err := e.uc.GetBalance(eCtx.Request().Context(), usr)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusNotFound,
			fmt.Sprintf("err in ServerHandler.GetBalance(): %v", err.Error()))
	}

	resp := &userBalance{
		UserId: usr.UserId,
		Ub: balance{
			CurAmount:   bal.Balance,
			AvailAmount: bal.Balance,
		},
	}

	return eCtx.JSON(http.StatusOK, resp)
}

func (e *ServerHandler) AddFunds(eCtx echo.Context) error {
	balanceOp, err := e.parseBalanceOperationBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.AddFundsToUser.parseBalanceOperationBody(): %v", err))
	}

	return e.uc.AddFundsToUser(eCtx.Request().Context(), balanceOp)
}
