package postgres

import (
	"context"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
	"money-manager/money-manager/repository"
)

const (
	GetSumPerServiceInYearMonthSqlCmd = `select min(r.service_id) as s_id, sum(r.amount) from report r 
											   where date_part('year', created) = $1 
											   and date_part('month', created) = $2
											   group by r.service_id order by s_id`
)

func (e *pgMoneyManagerRepo) MakeReportMoneyPerService(ctx context.Context, reportPeriod entity.ReportPeriod) (entity.ReportMoneyPerService, error) {
	rows, err := e.db.Query(ctx, GetSumPerServiceInYearMonthSqlCmd, reportPeriod.Year, reportPeriod.Month)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}
	}

	resSet := make(entity.ReportMoneyPerService, 0)
	for rows.Next() {
		var r entity.ReportMoneyPerServiceRaw

		err := rows.Scan(&r.ServiceId, &r.Sum)
		if err != nil {
			return nil, errors.Wrap(err, "Err in: pgMoneyManagerRepo.MakeReportPerService.Scan()")
		}

		resSet = append(resSet, r)
	}

	return resSet, errors.Wrap(rows.Err(), "Err in: pgMoneyManagerRepo.MakeReportPerService")
}
