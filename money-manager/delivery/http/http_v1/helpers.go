package http_v1

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/entity"
)

func (e *ServerHandler) parseBalanceOperationBody(eCtx echo.Context) (entity.BalanceOperation, error) {
	contentTypes := eCtx.Request().Header.Get(echo.HeaderContentType)

	balanceOp := &entity.BalanceOperation{}

	if contentTypes != "" {
		for _, ct := range strings.Split(contentTypes, ";") {
			if strings.TrimSpace(ct) == echo.MIMEApplicationJSON {

				err := eCtx.Bind(balanceOp)
				if err != nil {
					return *balanceOp, err
				}

				return *balanceOp, nil
			}
		}
	}

	return *balanceOp, errors.New("Content-Type header is missing")
}

func (e *ServerHandler) isValidUserId(id string) bool {
	return id != ""
}

func (e *ServerHandler) noContentErrResponse(eCtx echo.Context, statusCode int, errMsg string) error {
	log.Error().Msg(errMsg)

	return eCtx.NoContent(statusCode)
}
