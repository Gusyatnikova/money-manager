package http_v1

import (
	"os"

	"github.com/labstack/echo/v4"

	"money-manager/money-manager/delivery"
)

const (
	YearParamName  ParamName = "y"
	MonthParamName ParamName = "m"
)

func (e *ServerHandler) MakeReportMoneyPerService(eCtx echo.Context) error {
	queryParams := eCtx.QueryParams()

	year := queryParams.Get(string(YearParamName))
	month := queryParams.Get(string(MonthParamName))

	if year == "" || month == "" {
		return delivery.ErrBadRequestParams
	}

	filepath, filename, err := e.uc.MakeReportMoneyPerService(eCtx.Request().Context(), year, month)
	if err != nil {
		return err
	}

	defer os.Remove(filepath)

	return eCtx.Attachment(filepath, filename)
}
