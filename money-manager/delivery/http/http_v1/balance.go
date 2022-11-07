package http_v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"money-manager/money-manager/delivery"
	"money-manager/money-manager/entity"
)

type balance struct {
	Current   entity.MoneyAmount `json:"current_amount"`
	Available entity.MoneyAmount `json:"available_amount"`
	Unit      string             `json:"unit"`
}

type userBalanceResp struct {
	UserId string  `json:"user_id"`
	Ub     balance `json:"amount"`
}

type money struct {
	Value string `json:"amount"`
	Unit  string `json:"unit"`
}

type userMoneyReqBody struct {
	UserId string `json:"user_id"`
	Money  money  `json:"money"`
}

type transferMoneyReqBody struct {
	FromUserId string `json:"from_user_id"`
	ToUserId   string `json:"to_user_id"`
	Money      money  `json:"money"`
}

type transferMoneyResBody struct {
	FromUser userBalanceResp `json:"sender"`
	ToUser   userBalanceResp `json:"recipient"`
}

type ParamName string

const UserIdParamName ParamName = "user_id"

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	usr := eCtx.QueryParam(string(UserIdParamName))

	bal, err := e.uc.GetBalance(eCtx.Request().Context(), entity.UserId(usr))
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, makeUserBalanceResponse(entity.UserId(usr), bal))
}

func (e *ServerHandler) AddMoney(eCtx echo.Context) error {
	reqBody, err := parseUserMoneyBody(eCtx)
	if err != nil {
		return err
	}

	err = e.uc.AddMoneyToUser(eCtx.Request().Context(), entity.UserId(reqBody.UserId), userMoneyReqBodyToMoney(reqBody))
	if err != nil {
		return err
	}

	eCtx.QueryParams().Set(string(UserIdParamName), reqBody.UserId)

	return e.GetBalance(eCtx)
}

func (e *ServerHandler) DebitMoney(eCtx echo.Context) error {
	reqBody, err := parseUserMoneyBody(eCtx)
	if err != nil {
		return err
	}

	err = e.uc.DebitMoney(eCtx.Request().Context(), entity.UserId(reqBody.UserId), userMoneyReqBodyToMoney(reqBody))
	if err != nil {
		return err
	}

	eCtx.QueryParams().Set(string(UserIdParamName), reqBody.UserId)

	return e.GetBalance(eCtx)
}

func (e *ServerHandler) TransferMoney(eCtx echo.Context) error {
	reqBody, err := parseUserTransferReqBody(eCtx)
	if err != nil {
		return err
	}

	err = e.uc.TransferMoneyUserToUser(
		eCtx.Request().Context(),
		entity.UserId(reqBody.FromUserId),
		entity.UserId(reqBody.ToUserId),
		transferMoneyReqBodyToMoney(reqBody),
	)
	if err != nil {
		return err
	}

	userIdFrom := entity.UserId(reqBody.FromUserId)
	userIdTo := entity.UserId(reqBody.ToUserId)

	userBalFrom, err := e.uc.GetBalance(eCtx.Request().Context(), userIdFrom)
	if err != nil {
		return err
	}
	userBalTo, err := e.uc.GetBalance(eCtx.Request().Context(), userIdTo)
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, makeTransferMoneyResBody(userIdFrom, userIdTo, userBalFrom, userBalTo))
}

func parseUserMoneyBody(eCtx echo.Context) (userMoneyReqBody, error) {
	frBody := userMoneyReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return frBody, delivery.ErrBadContentType
	}

	err := eCtx.Bind(&frBody)
	if err != nil {
		return frBody, delivery.ErrBadRequestBody
	}

	return frBody, nil
}

func parseUserTransferReqBody(eCtx echo.Context) (transferMoneyReqBody, error) {
	frBody := transferMoneyReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return frBody, delivery.ErrBadContentType
	}

	err := eCtx.Bind(&frBody)
	if err != nil {
		return frBody, delivery.ErrBadRequestBody
	}

	return frBody, nil
}

func makeUserBalanceResponse(usr entity.UserId, bal entity.Balance) userBalanceResp {
	return userBalanceResp{
		UserId: string(usr),
		Ub: balance{
			Current:   bal.Current,
			Available: bal.Available,
			Unit:      "kop",
		},
	}
}

func makeTransferMoneyResBody(usrFrom entity.UserId, usrTo entity.UserId, balFrom entity.Balance, balTo entity.Balance) transferMoneyResBody {
	return transferMoneyResBody{
		FromUser: makeUserBalanceResponse(usrFrom, balFrom),
		ToUser:   makeUserBalanceResponse(usrTo, balTo),
	}
}

func userMoneyReqBodyToMoney(req userMoneyReqBody) entity.Money {
	return entity.Money{
		Value: req.Money.Value,
		Unit:  req.Money.Unit,
	}
}

func transferMoneyReqBodyToMoney(req transferMoneyReqBody) entity.Money {
	return entity.Money{
		Value: req.Money.Value,
		Unit:  req.Money.Unit,
	}
}
