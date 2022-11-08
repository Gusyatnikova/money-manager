package http_v1

import (
	"os"

	"github.com/labstack/echo/v4"

	"money-manager/money-manager/delivery"
)

type reportRequest struct {
	Year  string `query:"y"`
	Month string `query:"m"`
	Type  string `query:"type"`
}

func (e *ServerHandler) MakeReport(eCtx echo.Context) error {
	reportReq := reportRequest{}

	err := (&echo.DefaultBinder{}).BindQueryParams(eCtx, &reportReq)
	if err != nil {
		return delivery.ErrBadRequestParams
	}

	filepath, filename, err := e.uc.MakeReport(
		eCtx.Request().Context(), reportReq.Type, reportReq.Year, reportReq.Month)
	if err != nil {
		return err
	}

	defer os.Remove(filepath)

	return eCtx.Attachment(filepath, filename)
}
