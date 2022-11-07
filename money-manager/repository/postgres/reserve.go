package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
	"money-manager/money-manager/repository"
)

const (
	ReserveMoneySqlCmd = `insert into public.reserve (user_id , service_id, order_id, amount)
                          values ($1, $2, $3, $4)`

	DebitMoneyFromReserveSqlCmd = `update public.reserve set amount = amount - $4 
								   where user_id = $1 and service_id = $2 and order_id = $3 returning amount`

	AddMoneyToUserFromReserveSqlCmd = `update public.user as u 
                                       set (amount, updated) = (u.amount + r.amount, now()) 
									   from (select amount from reserve where user_id = $1 and service_id = $2 and order_id = $3) as r 
									   where id = $1`

	DeleteReserveSqlCmd = `delete from public.reserve 
						   where user_id = $1 and service_id = $2 and order_id = $3`

	GetAmountFromReserveSqlCmd = `select amount from public.reserve 
								  where user_id = $1 and service_id = $2 and order_id = $3`

	AddReserveToReportSqlCmd = `insert into public.report values($1, $2, $3, now())`
)

func (e *pgMoneyManagerRepo) ReserveMoney(ctx context.Context, res entity.Reserve) error {
	reserveMoneyFn := func(tx pgx.Tx) error {
		//subtract amount from public.user
		if _, err := tx.Exec(ctx, DebitMoneyFromUserSqlCmd, res.UserId, res.MoneyAmount); err != nil {
			return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Exec() for DebitMoneyFromUserSqlCmd")
		}

		//add amount to public.reserve
		_, err := tx.Exec(ctx, ReserveMoneySqlCmd, res.UserId, res.ServiceId, res.OrderId, res.MoneyAmount)
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.ReserveMoney.Exec() for ReserveMoneySqlCmd")
	}

	err := e.RunTx(ctx, reserveMoneyFn)

	return errors.Wrap(err, "Err in: pgMoneyManagerRepo.ReserveMoney.RunTx()")
}

func (e *pgMoneyManagerRepo) GetReserve(ctx context.Context, res entity.Reserve) (entity.Reserve, error) {
	var (
		resAmount uint64
		newRes    entity.Reserve
	)
	row := e.db.QueryRow(ctx, GetAmountFromReserveSqlCmd, res.UserId, res.ServiceId, res.OrderId)

	err := row.Scan(&resAmount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newRes, repository.ErrNotFound
		}
	}

	newRes.MoneyAmount = entity.MoneyAmount(resAmount)
	return newRes, errors.Wrap(err, "Err in: pgMoneyManagerRepo.GetReserve.Scan(): ")
}

func (e *pgMoneyManagerRepo) AcceptReserve(ctx context.Context, res entity.Reserve) error {
	acceptReserveFn := func(tx pgx.Tx) error {
		//subtract amount from public.reserve
		newAmount := uint64(0)

		err := tx.QueryRow(ctx, DebitMoneyFromReserveSqlCmd, res.UserId, res.ServiceId, res.OrderId, res.MoneyAmount).Scan(&newAmount)
		if err != nil {
			return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AcceptReserve.Exec() for DebitMoneyFromReserveSqlCmd")
		}

		//if reserve still have money return it to user
		if newAmount > 0 {
			if _, err := tx.Exec(ctx, AddMoneyToUserSqlCmd, res.UserId, newAmount); err != nil {
				return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AcceptReserve.Exec() for AddMoneyToUserFromReserveSqlCmd")
			}
		}

		//delete reserve
		if _, err := tx.Exec(ctx, DeleteReserveSqlCmd, res.UserId, res.ServiceId, res.OrderId); err != nil {
			return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AcceptReserve.Exec() for DeleteReserveSqlCmd")
		}

		//add accepted reserve to report
		_, err = tx.Exec(ctx, AddReserveToReportSqlCmd, ulid.Make().String(), res.ServiceId, res.MoneyAmount)

		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AcceptReserve")
	}

	err := e.RunTx(ctx, acceptReserveFn)

	return errors.Wrap(err, "Err in: pgMoneyManagerRepo.ReserveMoney.RunTx()")
}

func (e *pgMoneyManagerRepo) CancelReserve(ctx context.Context, res entity.Reserve) error {
	cancelReserveFn := func(tx pgx.Tx) error {
		//add money from reserve to user
		_, err := tx.Exec(ctx, AddMoneyToUserFromReserveSqlCmd, res.UserId, res.ServiceId, res.OrderId)
		if err != nil {
			return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.CancelReserve.Exec() for AddMoneyToUserFromReserveSqlCmd")
		}

		//delete reserve
		if _, err := tx.Exec(ctx, DeleteReserveSqlCmd, res.UserId, res.ServiceId, res.OrderId); err != nil {
			return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.CancelReserve.Exec() for DeleteReserveSqlCmd")
		}

		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.CancelReserve")
	}

	err := e.RunTx(ctx, cancelReserveFn)

	return errors.Wrap(err, "Err in: pgMoneyManagerRepo.CancelReserve.RunTx()")
}
