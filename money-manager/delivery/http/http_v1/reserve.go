package http_v1

import (
	"github.com/labstack/echo/v4"
	"money-manager/money-manager/entity"
)

type addReserveReqBody struct {
	UserId    string `json:"user_id"`
	ServiceId string `json:"service_id"`
	OrderId   string `json:"order_id"`
	Money     money  `json:"funds"`
}

func (e *ServerHandler) AddReserve(eCtx echo.Context) error {
	reqBody, err := parseAddReserveReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	return e.uc.Reserve(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reqBody.Money.Value, reqBody.Money.Unit)
}

func (e *ServerHandler) RevokeReserve(eCtx echo.Context) error {
	reqBody, err := parseAddReserveReqBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, err)
	}

	return e.uc.RevokeReserve(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reqBody.Money.Value, reqBody.Money.Unit)
}

func (e *ServerHandler) AcceptReserve(eCtx echo.Context) error {
	return nil
}

func parseAddReserveReqBody(eCtx echo.Context) (addReserveReqBody, error) {
	addReserveBody := addReserveReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return addReserveBody, ErrBadContentType
	}

	err := eCtx.Bind(&addReserveBody)
	if err != nil {
		return addReserveBody, ErrBadRequestBody
	}

	return addReserveBody, nil
}

func reserveReqBodyToReserve(res addReserveReqBody) entity.Reserve {
	return entity.Reserve{
		UserId:    entity.UserId(res.UserId),
		ServiceId: entity.ServiceId(res.ServiceId),
		OrderId:   entity.OrderId(res.OrderId),
	}
}
