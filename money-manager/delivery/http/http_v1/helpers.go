package http_v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/delivery"
	"money-manager/money-manager/usecase"
)

//noContentErrResponse response with no content and with http status code selected based on err
func (e *ServerHandler) noContentErrResponse(eCtx echo.Context, err error) error {
	log.Error().Msg(errors.Wrap(err, "Err in ServerHandler: ").Error())

	return eCtx.NoContent(getErrStatusCode(err))
}

//getErrStatusCode returns http status code based on getted err
func getErrStatusCode(err error) int {
	if errors.Is(err, usecase.ErrInvalidUser) {
		return http.StatusBadRequest
	}
	if errors.Is(err, usecase.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, usecase.ErrInvalidMoney) {
		return http.StatusBadRequest
	}
	if errors.Is(err, usecase.ErrMoneyLimitIsExceeded) {
		return http.StatusBadRequest
	}
	if errors.Is(err, usecase.ErrNotEnoughMoney) {
		return http.StatusBadRequest
	}

	if errors.Is(err, usecase.ErrInvalidReserve) {
		return http.StatusBadRequest
	}

	if errors.Is(err, delivery.ErrBadRequestBody) {
		return http.StatusBadRequest
	}
	if errors.Is(err, delivery.ErrBadContentType) {
		return http.StatusUnsupportedMediaType
	}

	return http.StatusInternalServerError
}

//isRequestBodyIsJSON returns true if header Contain-type with value application/json are in the request
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
