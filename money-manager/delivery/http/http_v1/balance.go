package http_v1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (e *ServerHandler) GetBalance(eCtx echo.Context) error {
	eCtx.NoContent(http.StatusOK)

	return nil
}

func (e *ServerHandler) AddFunds(eCtx echo.Context) error {
	balanceOp, err := e.parseBalanceOperationBody(eCtx)
	if err != nil {
		return e.noContentErrResponse(eCtx, http.StatusBadRequest,
			fmt.Sprintf("err in ServerHandler.AddFunds.parseBalanceOperationBody(): %v", err))
	}

	return e.uc.AddFunds(eCtx.Request().Context(), balanceOp)
}
