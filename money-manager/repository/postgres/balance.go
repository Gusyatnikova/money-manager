package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
	"money-manager/money-manager/usecase"
)

type pgMoneyManagerRepo struct {
	db *pgxpool.Pool
}

func NewPgMoneyManagerRepo(db *pgxpool.Pool) usecase.MoneyManagerRepo {
	return &pgMoneyManagerRepo{
		db: db,
	}
}

func (e *pgMoneyManagerRepo) AddFundsToUser(ctx context.Context, usr entity.User, fnd entity.Fund) error {
	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Begin()")
	}

	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := `
		update public.user set (amount, updated) = (public.user.amount + $2, now())
		where user_id = $1
	`

	stmt, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Prepare()")
	}

	if _, err := ts.Exec(ctx, stmt.SQL, usr.UserId, fnd.Amount); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Exec()")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) AddUser(ctx context.Context, usr entity.User, fnd entity.Fund) error {
	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Begin()")
	}

	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := `
		insert into public.user values(default, $1, $2, now(), null)
	`

	stmt, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Prepare()")
	}

	if _, err := ts.Exec(ctx, stmt.SQL, usr.UserId, fnd.Amount); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Exec()")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.AddFundsToUser.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) GetBalance(ctx context.Context, usr entity.User) (entity.Balance, error) {
	ub := entity.Balance{}

	sqlCmd := `
		select amount from public.user where user_id = $1
    `

	row := e.db.QueryRow(ctx, sqlCmd, usr.UserId)
	var amount uint64
	err := row.Scan(&amount)
	if err != nil {
		return ub, errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.GetBalance.Scan()")
	}

	ub.Available = entity.Fund{Amount: amount}
	ub.Current = entity.Fund{Amount: amount}

	return ub, nil
}
