package http_v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/delivery"
	"money-manager/money-manager/usecase"
)

//ErrorHandlerMiddleware handle responses which ended with not nil error
var ErrorHandlerMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		if err := next(eCtx); err != nil {
			return getHttpErr(err)
		}
		return nil
	}
}

//errResponse log error and return new http error with custom status code and message
func getHttpErr(err error) error {
	log.Error().Msg(errors.Wrap(err, "Err in ServerHandler: ").Error())

	return echo.NewHTTPError(errToHttpErr(err))
}

//getErrStatusCode returns http status code based on err and appropriate error message
func errToHttpErr(err error) (int, string) {
	//usecase common errors
	if errors.Is(err, usecase.ErrInvalidUser) {
		return http.StatusBadRequest, usecase.ErrInvalidUser.Error()
	}
	if errors.Is(err, usecase.ErrNotFound) {
		return http.StatusNotFound, usecase.ErrNotFound.Error()
	}
	if errors.Is(err, usecase.ErrInvalidMoney) {
		return http.StatusBadRequest, usecase.ErrInvalidMoney.Error()
	}
	if errors.Is(err, usecase.ErrDuplication) {
		return http.StatusConflict, usecase.ErrDuplication.Error()
	}

	//usecase balance errors
	if errors.Is(err, usecase.ErrMoneyLimitIsExceeded) {
		return http.StatusBadRequest, usecase.ErrMoneyLimitIsExceeded.Error()
	}
	if errors.Is(err, usecase.ErrNotEnoughMoney) {
		return http.StatusBadRequest, usecase.ErrNotEnoughMoney.Error()
	}

	//usecase reserve errors
	if errors.Is(err, usecase.ErrInvalidReserve) {
		return http.StatusBadRequest, usecase.ErrInvalidReserve.Error()
	}

	//delivery errors
	if errors.Is(err, delivery.ErrBadRequestBody) {
		return http.StatusBadRequest, delivery.ErrBadRequestBody.Error()
	}
	if errors.Is(err, delivery.ErrBadContentType) {
		return http.StatusUnsupportedMediaType, delivery.ErrBadContentType.Error()
	}
	if errors.Is(err, delivery.ErrBadRequestParams) {
		return http.StatusBadRequest, delivery.ErrBadRequestParams.Error()
	}

	//report errors
	if errors.Is(err, usecase.ErrInvalidReportInfo) {
		return http.StatusBadRequest, usecase.ErrInvalidReportInfo.Error()
	}

	return http.StatusInternalServerError, usecase.ErrInternalError.Error()
}
