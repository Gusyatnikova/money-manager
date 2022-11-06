package http_v1

import (
	"github.com/labstack/echo/v4"
)

func (e *ServerHandler) HealthCheck(eCtx echo.Context) error {
	return e.uc.HealthCheck(eCtx.Request().Context())
}
