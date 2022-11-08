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
} //@name ReportRequest

// MakeReport godoc
// @Summary     Return attachment with generated report for required report type, year and month
// @Description type=money_per_service for money per Service report
// @Tags        Report operations
// @Param       type query  string true "type of report"
// @Param       m    query  string true "month for data selection"
// @Param       y    query  string true "year for data selection"
// @Header      200  header string "Content-Disposition header with file name"
// @Success     200  string ContentDisposition
// @Failure     400  "Invalid request for report"
// @Router      /reports [get]
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
