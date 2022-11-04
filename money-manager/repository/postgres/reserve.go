package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

const (
	AddFundsToReserveSqlCmd = `insert into public.reserve (user_id , service_id, order_id, amount)
                               	values (
									(select id from public.user as u
									where u.user_id = $1
									), $2, $3, $4
								)`
	AcceptReserveFundsSqlCmd = `rec`
)

func (e *pgMoneyManagerRepo) AddReserve(ctx context.Context, res entity.Reserve) error {
	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddReserve.Begin()")
	}

	//subtract amount from public.user
	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := DebitFundsFromUserSqlCmd

	stmtUser, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddReserve.Prepare() for DebitFundsFromUserSqlCmd")
	}

	if _, err := ts.Exec(ctx, stmtUser.SQL, res.UserId, res.Amount); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddReserve.Exec() for DebitFundsFromUserSqlCmd")
	}

	//add amount to public.reserve
	prepPart = fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd = AddFundsToReserveSqlCmd

	stmtReserve, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddReserve.Prepare() for AddFundsToReserveSqlCmd")
	}

	if _, err := ts.Exec(ctx, stmtReserve.SQL, res.UserId, res.ServiceId, res.OrderId, res.Amount); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddReserve.Exec() for AddFundsToReserveSqlCmd")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddReserve.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) AcceptReserve(ctx context.Context, resId entity.ReserveId) error {
	return nil
}
