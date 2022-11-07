package http_v1

import (
	"github.com/labstack/echo/v4"

	"money-manager/money-manager/delivery"
	"money-manager/money-manager/entity"
)

type ReserveKey struct {
	UserId    string `json:"user_id"`
	ServiceId string `json:"service_id"`
	OrderId   string `json:"order_id"`
}

type reserveReqBody struct {
	ReserveKey ReserveKey `json:"reserve_id"`
	Money      money      `json:"money"`
}

func (e *ServerHandler) AddReserve(eCtx echo.Context) error {
	reqBody, err := parseReserveReqBody(eCtx)
	if err != nil {
		return err
	}

	err = e.uc.ReserveMoney(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reserveReqBodyToMoney(reqBody))
	if err != nil {
		return err
	}

	eCtx.QueryParams().Set(string(UserIdParamName), reqBody.ReserveKey.UserId)

	return e.GetBalance(eCtx)
}

func (e *ServerHandler) CancelReserve(eCtx echo.Context) error {
	reqBody, err := parseCancelReserveReqBody(eCtx)
	if err != nil {
		return err
	}

	reserveReqBody := reserveReqBody{
		ReserveKey: reqBody,
		Money:      money{},
	}

	err = e.uc.CancelReserve(eCtx.Request().Context(), reserveReqBodyToReserve(reserveReqBody))
	if err != nil {
		return err
	}

	return nil
}

func (e *ServerHandler) AcceptReserve(eCtx echo.Context) error {
	reqBody, err := parseReserveReqBody(eCtx)
	if err != nil {
		return err
	}

	err = e.uc.AcceptReserve(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reserveReqBodyToMoney(reqBody))
	if err != nil {
		return err
	}

	eCtx.QueryParams().Set(string(UserIdParamName), reqBody.ReserveKey.UserId)

	return e.GetBalance(eCtx)
}

func parseReserveReqBody(eCtx echo.Context) (reserveReqBody, error) {
	addReserveBody := reserveReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return addReserveBody, delivery.ErrBadContentType
	}

	err := eCtx.Bind(&addReserveBody)
	if err != nil {
		return addReserveBody, delivery.ErrBadRequestBody
	}

	return addReserveBody, nil
}

func parseCancelReserveReqBody(eCtx echo.Context) (ReserveKey, error) {
	cancelReserveBody := ReserveKey{}

	if !isRequestBodyIsJSON(eCtx) {
		return cancelReserveBody, delivery.ErrBadContentType
	}

	err := eCtx.Bind(&cancelReserveBody)
	if err != nil {
		return cancelReserveBody, delivery.ErrBadRequestBody
	}

	return cancelReserveBody, nil
}

func reserveReqBodyToReserve(res reserveReqBody) entity.Reserve {
	return entity.Reserve{
		UserId:    entity.UserId(res.ReserveKey.UserId),
		ServiceId: entity.ServiceId(res.ReserveKey.ServiceId),
		OrderId:   entity.OrderId(res.ReserveKey.OrderId),
	}
}

func reserveReqBodyToMoney(res reserveReqBody) entity.Money {
	return entity.Money{
		Value: res.Money.Value,
		Unit:  res.Money.Unit,
	}
}
