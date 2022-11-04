package http_v1

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"strings"
)

func noContentErrResponse(eCtx echo.Context, statusCode int, errMsg string) error {
	log.Error().Msg(errMsg)

	return eCtx.NoContent(statusCode)
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
