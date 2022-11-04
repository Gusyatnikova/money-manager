package http_v1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

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
		return noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.AddReserve.parseAddReserveReqBody(): %v", err))
	}

	return e.uc.Reserve(eCtx.Request().Context(), reserveReqBodyToReserve(reqBody), reqBody.Money.Value, reqBody.Money.Unit)
}

func (e *ServerHandler) AcceptReserve(eCtx echo.Context) error {
	return nil
}

func parseAddReserveReqBody(eCtx echo.Context) (addReserveReqBody, error) {
	addReserveBody := addReserveReqBody{}

	if !isRequestBodyIsJSON(eCtx) {
		return addReserveBody, errors.New("Content-Type application/json is missing")
	}

	err := eCtx.Bind(&addReserveBody)
	if err != nil {
		return addReserveBody, errors.Wrap(err, "Unable parse request body")
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
