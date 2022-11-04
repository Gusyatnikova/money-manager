package http_v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"

	"money-manager/money-manager/entity"
)

type balance struct {
	CurAmount   entity.Fund `json:"current_amount"`
	AvailAmount entity.Fund `json:"available_amount"`
	Unit        string      `json:"unit"`
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

type transferFundsReqBody struct {
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	Amount     string `json:"amount"`
	Unit       string `json:"unit"`
}

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	usr := eCtx.QueryParam("user_id")

	bal, err := e.uc.GetBalance(eCtx.Request().Context(), entity.UserId(usr))
	if err != nil {
		return noContentErrResponse(eCtx, http.StatusNotFound,
			fmt.Sprintf("err in ServerHandler.GetBalance(): %v", err.Error()))
	}

	return eCtx.JSON(http.StatusOK, makeUserBalanceResponse(entity.UserId(usr), bal))
}

func (e *ServerHandler) AddFunds(eCtx echo.Context) error {
	reqBody, err := parseUserAmountBody(eCtx)
	if err != nil {
		return noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.AddFundsToUser.parseUserAmountBody(): %v", err))
	}

	return e.uc.AddFundsToUser(eCtx.Request().Context(), entity.UserId(reqBody.UserId), reqBody.Amount, reqBody.Unit)
}

func (e *ServerHandler) DebitFunds(eCtx echo.Context) error {
	reqBody, err := parseUserAmountBody(eCtx)
	if err != nil {
		return noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.DebitFunds.parseUserAmountBody(): %v", err))
	}

	return e.uc.DebitFunds(eCtx.Request().Context(), entity.UserId(reqBody.UserId), reqBody.Amount, reqBody.Unit)
}

func (e *ServerHandler) TransferFunds(eCtx echo.Context) error {
	reqBody, err := parseUserTransferReqBody(eCtx)
	if err != nil {
		return noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.TransferFunds.parseUserTransferReqBody(): %v", err))
	}

	return e.uc.TransferFundsUserToUser(
		eCtx.Request().Context(),
		entity.UserId(reqBody.FromUserId),
		entity.UserId(reqBody.ToUserId),
		reqBody.Amount,
		reqBody.Unit)
}

func parseUserAmountBody(eCtx echo.Context) (fundsReqBody, error) {
	frBody := fundsReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return frBody, errors.New("Content-Type application/json is missing")
	}

	err := eCtx.Bind(&frBody)
	if err != nil {
		return frBody, errors.Wrap(err, "Unable parse request body")
	}

	return frBody, nil
}

func parseUserTransferReqBody(eCtx echo.Context) (transferFundsReqBody, error) {
	frBody := transferFundsReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return frBody, errors.New("Content-Type application/json is missing")
	}

	err := eCtx.Bind(&frBody)
	if err != nil {
		return frBody, errors.Wrap(err, "Unable parse request body")
	}

	return frBody, nil
}

func makeUserBalanceResponse(usr entity.UserId, bal entity.Balance) userBalanceResp {
	return userBalanceResp{
		UserId: string(usr),
		Ub: balance{
			CurAmount:   bal.Current,
			AvailAmount: bal.Available,
			Unit:        "kop",
		},
	}
}
