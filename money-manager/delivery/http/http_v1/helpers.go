package http_v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/usecase/money_manager"
)

func (e *ServerHandler) noContentErrResponse(eCtx echo.Context, err error) error {
	log.Error().Msg(errors.Wrap(err, "Err in ServerHandler: ").Error())

	return eCtx.NoContent(getErrStatusCode(err))
}

func getErrStatusCode(err error) int {
	switch err {
	case money_manager.ErrInvalidUser:
		return http.StatusBadRequest
	case money_manager.ErrNotFound:
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func isRequestBodyIsJSON(eCtx echo.Context) bool {
	contentTypes := eCtx.Request().Header.Get(echo.HeaderContentType)

	if contentTypes != "" {
		for _, ct := range strings.Split(contentTypes, ";") {
			if strings.TrimSpace(ct) == echo.MIMEApplicationJSON {
				return true
			}
		}
	}

	return false
}
