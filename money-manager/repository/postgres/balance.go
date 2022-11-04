package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

const (
	AddUserSqlCmd        = `insert into public.user values($1, $2, $3, now(), null)`
	AddFundsToUserSqlCmd = `update public.user set (amount, updated) = (public.user.amount + $2, now())
								where user_id = $1`
	GetUserBalanceSqlCmd = `select min(u.amount), coalesce(sum(r.amount), 0) from public.user u
								left join public.reserve r on u.id = r.user_id
							where u.user_id = $1 group by u.id`
)

func (e *pgMoneyManagerRepo) AddFundsToUser(ctx context.Context, usr entity.UserId, fnd entity.Fund) error {
	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Begin()")
	}

	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := AddFundsToUserSqlCmd

	stmt, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Prepare()")
	}

	if _, err := ts.Exec(ctx, stmt.SQL, usr, fnd); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Exec()")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) AddUser(ctx context.Context, usr entity.UserId, fnd entity.Fund) error {
	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Begin()")
	}

	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := AddUserSqlCmd

	stmt, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Prepare()")
	}

	if _, err := ts.Exec(ctx, stmt.SQL, generateID(), usr, fnd); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Exec()")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) GetBalance(ctx context.Context, usr entity.UserId) (entity.Balance, error) {
	sqlCmd := GetUserBalanceSqlCmd

	row := e.db.QueryRow(ctx, sqlCmd, usr)
	bal, err := e.ScanBalance(row)
	if err != nil {
		return bal, errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.GetBalance.Scan()")
	}

	return bal, nil
}

func (e *pgMoneyManagerRepo) ScanBalance(row pgx.Row) (entity.Balance, error) {
	bal := entity.Balance{}

	availVal := uint64(0)
	resVal := uint64(0)

	err := row.Scan(&availVal, &resVal)
	if err != nil {
		return bal, errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.GetBalance.Scan()")
	}

	bal.Available = entity.Fund(availVal)
	bal.Current = entity.Fund(availVal)

	if resVal != 0 {
		bal.Current = entity.Fund(availVal + resVal)
	}

	return bal, nil
}
