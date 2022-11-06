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
	Money      money      `json:"funds"`
}

func (e *ServerHandler) AddReserve(eCtx echo.Context) error {
	reqBody, err := parseReserveReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	err = e.uc.ReserveMoney(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reqBody.Money.Value, reqBody.Money.Unit)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	eCtx.QueryParams().Set(UserIdParamName, reqBody.ReserveKey.UserId)

	//todo: get this data from e.uc.AddMoneyToUser return
	return e.GetBalance(eCtx)
}

func (e *ServerHandler) CancelReserve(eCtx echo.Context) error {
	reqBody, err := parseCancelReserveReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	reserveReqBody := reserveReqBody{
		ReserveKey: reqBody,
		Money:      money{},
	}

	err = e.uc.CancelReserve(eCtx.Request().Context(), reserveReqBodyToReserve(reserveReqBody))
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	return nil
}

func (e *ServerHandler) AcceptReserve(eCtx echo.Context) error {
	reqBody, err := parseReserveReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	err = e.uc.AcceptReserve(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reqBody.Money.Value, reqBody.Money.Unit)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	eCtx.QueryParams().Set(UserIdParamName, reqBody.ReserveKey.UserId)

	//todo: get this data from e.uc.AddMoneyToUser return
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
