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
	Unit        string `json:"unit"`
}

type userBalanceResp struct {
	UserId string  `json:"user_id"`
	Ub     balance `json:"balance"`
}

type fundsReqBody struct {
	UserId string `json:"user_id"`
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	usr := entity.User{
		UserId: eCtx.Param("user_id"),
	}

	bal, err := e.uc.GetBalance(eCtx.Request().Context(), usr)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusNotFound,
			fmt.Sprintf("err in ServerHandler.GetBalance(): %v", err.Error()))
	}

	return eCtx.JSON(http.StatusOK, makeUserBalanceResponse(usr, bal))
}

func (e *ServerHandler) AddFunds(eCtx echo.Context) error {
	reqBody, err := e.parseUserAmountBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.AddFundsToUser.parseUserAmountBody(): %v", err))
	}

	return e.uc.AddFundsToUser(eCtx.Request().Context(), reqBodyToUser(reqBody), reqBody.Amount, reqBody.Unit)
}
