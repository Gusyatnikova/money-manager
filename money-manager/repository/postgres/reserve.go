package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

const (
	AddMoneyToReserveSqlCmd = `insert into public.reserve (user_id , service_id, order_id, amount)
                               	values (
									(select id from public.user as u
									where u.user_id = $1
									), $2, $3, $4
								)`
	AcceptReserveMoneySqlCmd = `rec`
)

func (e *pgMoneyManagerRepo) ReserveMoney(ctx context.Context, res entity.Reserve) error {
	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Begin()")
	}

	//subtract amount from public.user
	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := DebitMoneyFromUserSqlCmd

	stmtUser, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Prepare() for DebitMoneyFromUserSqlCmd")
	}

	if _, err := ts.Exec(ctx, stmtUser.SQL, res.UserId, res.Amount); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Exec() for DebitMoneyFromUserSqlCmd")
	}

	//add amount to public.reserve
	prepPart = fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd = AddMoneyToReserveSqlCmd

	stmtReserve, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Prepare() for AddMoneyToReserveSqlCmd")
	}

	if _, err := ts.Exec(ctx, stmtReserve.SQL, res.UserId, res.ServiceId, res.OrderId, res.Amount); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Exec() for AddMoneyToReserveSqlCmd")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) RevokeReserve(ctx context.Context, res entity.Reserve) error {
	return nil
}

func (e *pgMoneyManagerRepo) AcceptReserve(ctx context.Context, resId entity.ReserveId) error {
	return nil
}
