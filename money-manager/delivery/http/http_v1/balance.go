package http_v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"

	"money-manager/money-manager/entity"
)

type balance struct {
	CurMoney   entity.Fund `json:"current_amount"`
	AvailMoney entity.Fund `json:"available_amount"`
	Unit       string      `json:"unit"`
}

type userBalanceResp struct {
	UserId string  `json:"user_id"`
	Ub     balance `json:"amount"`
}

type money struct {
	Value string `json:"amount"`
	Unit  string `json:"unit"`
}

type fundsReqBody struct {
	UserId string `json:"user_id"`
	Money  money  `json:"funds"`
}

type transferFundsReqBody struct {
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	Money      money  `json:"funds"`
}

var ErrBadContentType = errors.New("Content-Type application/json is missing")
var ErrBadRequestBody = errors.New("Request body is incorrect")

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	usr := eCtx.QueryParam("user_id")

	bal, err := e.uc.GetBalance(eCtx.Request().Context(), entity.UserId(usr))
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	return eCtx.JSON(http.StatusOK, makeUserBalanceResponse(entity.UserId(usr), bal))
}

func (e *ServerHandler) AddFunds(eCtx echo.Context) error {
	reqBody, err := parseUserMoneyBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	err = e.uc.AddFundsToUser(eCtx.Request().Context(), entity.UserId(reqBody.UserId), reqBody.Money.Value, reqBody.Money.Unit)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	return eCtx.NoContent(http.StatusCreated)
}

func (e *ServerHandler) DebitFunds(eCtx echo.Context) error {
	reqBody, err := parseUserMoneyBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.DebitFunds.parseUserMoneyBody(): %v", err))
	}

	return e.uc.DebitFunds(eCtx.Request().Context(), entity.UserId(reqBody.UserId), reqBody.Money.Value, reqBody.Money.Unit)
}

func (e *ServerHandler) TransferFunds(eCtx echo.Context) error {
	reqBody, err := parseUserTransferReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.TransferFunds.parseUserTransferReqBody(): %v", err))
	}

	return e.uc.TransferFundsUserToUser(
		eCtx.Request().Context(),
		entity.UserId(reqBody.FromUserId),
		entity.UserId(reqBody.ToUserId),
		reqBody.Money.Value,
		reqBody.Money.Unit)
}

func parseUserMoneyBody(eCtx echo.Context) (fundsReqBody, error) {
	frBody := fundsReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return frBody, ErrWrongContentType
	}

	err := eCtx.Bind(&frBody)
	if err != nil {
		return frBody, ErrBadRequestBody
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
			CurMoney:   bal.Current,
			AvailMoney: bal.Available,
			Unit:       "kop",
		},
	}
}
