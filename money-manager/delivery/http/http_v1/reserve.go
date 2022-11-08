package http_v1

import (
	"github.com/labstack/echo/v4"

	"money-manager/money-manager/delivery"
	"money-manager/money-manager/entity"
)

type ReserveKey struct {
	UserId    string `json:"user_id" example:"user123"`
	ServiceId string `json:"service_id" example:"delivery"`
	OrderId   string `json:"order_id" example:"123"`
} //@name ReserveKey

type reserveReqBody struct {
	ReserveKey ReserveKey `json:"reserve_id"`
	Money      money      `json:"money"`
} //@name ReserveRequestBody

// AddReserve godoc
// @Summary Reserve user's money to separate account
// @Tags    Reserve operations
// @Param   request body http_v1.reserveReqBody true "ID's of user, service, order and money to reserve"
// @Produce json
// @Success 200 {object} http_v1.userBalanceResp
// @Failure 400 "Invalid id's of reserve | Invalid money | Insufficient funds to complete the operation"
// @Failure 404 "Requested resource is not found"
// @Router  /reserves/add [post]
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

// CancelReserve godoc
// @Summary Return reserved money to user's account and delete reserve
// @Tags    Reserve operations
// @Param   request body ReserveKey true "ID's of user, service, order and money to reserve"
// @Produce json
// @Success 200 {object} userBalanceResp
// @Failure 400 "Invalid id's of reserve"
// @Failure 404 "Requested resource is not found"
// @Router  /reserves/cancel [delete]
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

	eCtx.QueryParams().Set(string(UserIdParamName), reqBody.UserId)

	return e.GetBalance(eCtx)
}

// AcceptReserve godoc
// @Summary     Debit money from reserve and add to report.
// @Description "If accepted money = reserved money, add operation to report and delete reserve, if accepted money < reserved money, add to user's balance (reserved money - accepted money) and add operation to report, then delete reserve. If accepted money > reserved money then error 400"
// @Tags        Reserve operations
// @Param       request body reserveReqBody true "ID's of user, service, order and money to reserve"
// @Produce     json
// @Success     200 {object} userBalanceResp
// @Failure     400 "Invalid id's of reserve or | Invalid money | Insufficient funds to complete the operation"
// @Failure     404 "Requested resource is not found"
// @Router      /reserves/accept [delete]
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
