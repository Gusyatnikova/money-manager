package money_manager

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"money-manager/money-manager/delivery"
	"money-manager/money-manager/entity"
	"money-manager/money-manager/usecase"
)

const moneyPerServiceReportType = "money_per_service"

type reportInfo struct {
	Type   string
	Period entity.ReportPeriod
}

//MakeReport make report based on reportType value and returns its full path and filename
func (e *moneyManager) MakeReport(ctx context.Context, reportType, yearStr, monthStr string) (string, string, error) {
	reportPeriod, err := yearMonthToReportPeriod(yearStr, monthStr)
	if err != nil {
		return "", "", usecase.ErrInvalidReportInfo
	}

	reportInfo := reportInfo{
		Type:   reportType,
		Period: reportPeriod,
	}

	switch strings.ToLower(reportType) {
	case moneyPerServiceReportType:
		return e.generateMoneyPerServiceReport(ctx, reportInfo)
	default:
		return "", "", usecase.ErrInvalidReportInfo
	}
}

//generateMoneyPerServiceReport generate report of moneyPerServiceReportType type
func (e *moneyManager) generateMoneyPerServiceReport(ctx context.Context, reportInfo reportInfo) (string, string, error) {
	reportData, err := e.repo.MakeReportMoneyPerService(ctx, reportInfo.Period)
	if err != nil {
		return "", "", errors.Wrap(err, "err in moneyManager.MakeReport().MakeReport():")
	}

	reportCSV, err := generateCSV(reportData, reportInfo)
	if err != nil {
		return "", "", errors.Wrap(err, "err in moneyManager.MakeReport().generateCSV():")
	}

	reportPath, err := getFullPath(reportCSV)

	return reportPath, generateMoneyPerServiceReportName(reportInfo.Type),
		errors.Wrap(err, "err in moneyManager.MakeReport.generateMoneyPerServiceReport():")
}

//generateCSV generate .csv from report and returns its full path
func generateCSV(reportData entity.ReportMoneyPerService, reportInfo reportInfo) (*os.File, error) {
	data := reportMoneyPerServiceToSlice(reportData)

	reportFile, err := os.Create(generateMoneyPerServiceReportName(reportInfo.Type))
	if err != nil {
		return nil, errors.Wrap(err, "Err in moneyManager.MakeReport().Open():")
	}
	defer reportFile.Close()

	csvWriter := csv.NewWriter(reportFile)
	csvWriter.Comma = ';'

	for _, record := range data {
		if err := csvWriter.Write(record); err != nil {
			return nil, errors.Wrap(err, "Err in generateCSV.Write():")
		}
	}

	csvWriter.Flush()
	return reportFile, errors.Wrap(err, "Err in generateCSV()")
}

func getFullPath(file *os.File) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "Err in getFullPath.Getwd():")
	}

	return filepath.Join(path, file.Name()), nil
}

func generateMoneyPerServiceReportName(reportType string) string {
	now := time.Now()

	//report filename template is moneyPerServiceReportName_YYYY-MM-DD.csv
	return fmt.Sprintf("%s_%d-%02d-%02d.csv",
		reportType, now.Year(), now.Month(), now.Day())
}

func reportMoneyPerServiceToSlice(report entity.ReportMoneyPerService) [][]string {
	resultSl := make([][]string, 0)

	for _, raw := range report {
		rawSl := make([]string, 2)
		rawSl[0] = string(raw.ServiceId)
		rawSl[1] = raw.Sum

		resultSl = append(resultSl, rawSl)
	}

	return resultSl
}

func yearMonthToReportPeriod(yearStr, monthStr string) (entity.ReportPeriod, error) {
	var reportPeriod entity.ReportPeriod

	year, err := stringToYear(yearStr)
	if err != nil {
		return reportPeriod, err
	}

	month, err := stringToMonth(monthStr)
	if err != nil {
		return reportPeriod, err
	}

	reportPeriod.Year = year
	reportPeriod.Month = time.Month(month)

	return reportPeriod, nil
}

func stringToYear(str string) (int, error) {
	year, err := strconv.Atoi(str)
	if err != nil {
		return 0, delivery.ErrBadRequestParams
	}

	if year <= 0 {
		return 0, delivery.ErrBadRequestParams
	}

	return year, nil
}

func stringToMonth(str string) (int, error) {
	month, err := strconv.Atoi(str)
	if err != nil {
		return 0, delivery.ErrBadRequestParams
	}

	if month < 1 || month > 12 {
		return 0, delivery.ErrBadRequestParams
	}

	return month, nil
}
