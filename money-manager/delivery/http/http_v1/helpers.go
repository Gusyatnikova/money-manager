package http_v1

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/entity"
)

func (e *ServerHandler) parseUserAmountBody(eCtx echo.Context) (fundsReqBody, error) {
	contentTypes := eCtx.Request().Header.Get(echo.HeaderContentType)

	frBody := &fundsReqBody{}

	if contentTypes != "" {
		for _, ct := range strings.Split(contentTypes, ";") {
			if strings.TrimSpace(ct) == echo.MIMEApplicationJSON {

				err := eCtx.Bind(frBody)
				if err != nil {
					return *frBody, err
				}

				return *frBody, nil
			}
		}
	}

	return *frBody, errors.New("Content-Type header is missing")
}

func reqBodyToUser(rb fundsReqBody) entity.User {
	return entity.User{
		UserId: rb.UserId,
	}
}

func makeUserBalanceResponse(usr entity.User, bal entity.Balance) userBalanceResp {
	return userBalanceResp{
		UserId: usr.UserId,
		Ub: balance{
			CurAmount:   bal.Current.Amount,
			AvailAmount: bal.Available.Amount,
			Unit:        "kop",
		},
	}
}

func (e *ServerHandler) noContentErrResponse(eCtx echo.Context, statusCode int, errMsg string) error {
	log.Error().Msg(errMsg)

	return eCtx.NoContent(statusCode)
}
