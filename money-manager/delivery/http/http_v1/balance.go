package http_v1

import (
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

type transferFundsResBody struct {
	FromUser userBalanceResp `json:"sender"`
	ToUser   userBalanceResp `json:"recipient"`
}

const UserIdParam = "user_id"

var ErrBadContentType = errors.New("Content-Type application/json is missing")
var ErrBadRequestBody = errors.New("Request body is incorrect")

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	usr := eCtx.QueryParam(UserIdParam)

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

	eCtx.QueryParams().Set(UserIdParam, reqBody.UserId)

	return e.GetBalance(eCtx)
}

func (e *ServerHandler) DebitFunds(eCtx echo.Context) error {
	reqBody, err := parseUserMoneyBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	err = e.uc.DebitFunds(eCtx.Request().Context(), entity.UserId(reqBody.UserId), reqBody.Money.Value, reqBody.Money.Unit)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	eCtx.QueryParams().Set(UserIdParam, reqBody.UserId)

	return e.GetBalance(eCtx)
}

func (e *ServerHandler) TransferFunds(eCtx echo.Context) error {
	reqBody, err := parseUserTransferReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	err = e.uc.TransferFundsUserToUser(
		eCtx.Request().Context(),
		entity.UserId(reqBody.FromUserId),
		entity.UserId(reqBody.ToUserId),
		reqBody.Money.Value,
		reqBody.Money.Unit)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	userIdFrom := entity.UserId(reqBody.FromUserId)
	userIdTo := entity.UserId(reqBody.ToUserId)

	userBalFrom, err := e.uc.GetBalance(eCtx.Request().Context(), userIdFrom)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}
	userBalTo, err := e.uc.GetBalance(eCtx.Request().Context(), userIdTo)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	return eCtx.JSON(http.StatusOK, makeTransferFundsResBody(userIdFrom, userIdTo, userBalFrom, userBalTo))
}

func parseUserMoneyBody(eCtx echo.Context) (fundsReqBody, error) {
	frBody := fundsReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return frBody, ErrBadContentType
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
		return frBody, ErrBadContentType
	}

	err := eCtx.Bind(&frBody)
	if err != nil {
		return frBody, ErrBadRequestBody
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

func makeTransferFundsResBody(usrFrom entity.UserId, usrTo entity.UserId, balFrom entity.Balance, balTo entity.Balance) transferFundsResBody {
	return transferFundsResBody{
		FromUser: makeUserBalanceResponse(usrFrom, balFrom),
		ToUser:   makeUserBalanceResponse(usrTo, balTo),
	}
}
