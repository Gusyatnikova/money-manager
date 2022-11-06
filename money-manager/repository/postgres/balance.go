package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
	"money-manager/money-manager/repository"
)

const (
	AddUserSqlCmd = `insert into public.user values($1, $2, now(), null) returning (id, amount)`

	AddMoneyToUserSqlCmd = `update public.user set (amount, updated) = (amount + $2, now())
							where id = $1`

	GetUserBalanceSqlCmd = `select min(u.amount), coalesce(sum(r.amount), 0) from public.user u
							left join reserve r on u.id = r.user_id
							where u.id = $1 group by u.id`

	DebitMoneyFromUserSqlCmd = `update public.user set (amount, updated) = (amount - $2, now()) where id = $1`
)

func (e *pgMoneyManagerRepo) AddMoneyToUser(ctx context.Context, usr entity.UserId, fnd entity.Fund) error {
	addMoneyFn := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, AddMoneyToUserSqlCmd, usr, fnd)
		return errors.Wrap(err, "Err in: pgMoneyManagerRepo.AddMoneyToUser.Exec()")
	}

	err := e.RunTx(ctx, addMoneyFn)

	return errors.Wrap(err, "Err in: pgMoneyManagerRepo.AddMoneyToUser.RunTx()")
}

func (e *pgMoneyManagerRepo) AddUser(ctx context.Context, usr entity.UserId, fnd entity.Fund) error {
	addUserFn := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, AddUserSqlCmd, usr, fnd)
		return errors.Wrap(err, "Err in: pgMoneyManagerRepo.AddUser.Exec()")
	}

	err := e.RunTx(ctx, addUserFn)

	return errors.Wrap(err, "Err in: pgMoneyManagerRepo.AddUser.RunTx()")
}

func (e *pgMoneyManagerRepo) GetBalance(ctx context.Context, usr entity.UserId) (entity.Balance, error) {
	row := e.db.QueryRow(ctx, GetUserBalanceSqlCmd, usr)
	bal, err := e.ScanBalance(row)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return entity.Balance{}, repository.ErrNotFound
		}
	}

	return bal, errors.Wrap(err, "Err in: pgMoneyManagerRepo.GetBalance.ScanBalance(): ")
}

func (e *pgMoneyManagerRepo) TransferMoneyUserToUser(ctx context.Context,
	usrFrom entity.UserId, usrTo entity.UserId, fnd entity.Fund) error {

	transferMoneyFn := func(tx pgx.Tx) error {
		//subtract amount from usrFrom
		if _, err := tx.Exec(ctx, DebitMoneyFromUserSqlCmd, usrFrom, fnd); err != nil {
			return errors.Wrap(err, "Err in: pgMoneyManagerRepo.TransferMoneyUserToUser.Exec() for DebitMoneyFromUserSqlCmd")
		}

		//add amount to usrTo
		_, err := tx.Exec(ctx, AddMoneyToUserSqlCmd, usrTo, fnd)

		return errors.Wrap(err, "Err in: pgMoneyManagerRepo.TransferMoneyUserToUser.Exec() for AddMoneyToUserSqlCmd")
	}

	if err := e.RunTx(ctx, transferMoneyFn); err != nil {
		return errors.Wrap(err, "Err in: pgMoneyManagerRepo.AddUser.RunTx()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) DebitMoney(ctx context.Context, usrId entity.UserId, fnd entity.Fund) error {
	debitMoneyFn := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, DebitMoneyFromUserSqlCmd, usrId, fnd)
		return errors.Wrap(err, "Err in: pgMoneyManagerRepo.DebitMoney.Exec()")
	}

	err := e.RunTx(ctx, debitMoneyFn)

	return errors.Wrap(err, "Err in: pgMoneyManagerRepo.DebitMoney.RunTx()")
}

func (e *pgMoneyManagerRepo) ScanBalance(row pgx.Row) (entity.Balance, error) {
	bal := entity.Balance{}

	availVal := uint64(0)
	resVal := uint64(0)

	err := row.Scan(&availVal, &resVal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return bal, repository.ErrNotFound
		}
		return bal, errors.Wrap(err, "Err in: pgMoneyManagerRepo.ScanBalance.Scan()")
	}

	bal.Available = entity.Fund(availVal)
	bal.Current = entity.Fund(availVal)

	if resVal != 0 {
		bal.Current = entity.Fund(availVal + resVal)
	}

	return bal, nil
}
