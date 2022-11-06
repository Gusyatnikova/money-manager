package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

//use on conflict update?
const (
	AddUserSqlCmd        = `insert into public.user values($1, $2, $3, now(), null)`
	AddFundsToUserSqlCmd = `update public.user set (amount, updated) = (public.user.amount + $2, now())
							where user_id = $1`
	GetUserBalanceSqlCmd = `select min(u.amount), coalesce(sum(r.amount), 0) from public.user u
							left join public.reserve r on u.id = r.user_id
							where u.user_id = $1 group by u.id`
	DebitFundsFromUserSqlCmd = `update public.user set amount = amount - $2 where user_id = $1`
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

var ErrNotFound = errors.New("Err in pgMoneyManagerRepo: requested resource is not found")

func (e *pgMoneyManagerRepo) GetBalance(ctx context.Context, usr entity.UserId) (entity.Balance, error) {
	sqlCmd := GetUserBalanceSqlCmd

	row := e.db.QueryRow(ctx, sqlCmd, usr)
	bal, err := e.ScanBalance(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entity.Balance{}, ErrNotFound
		}
		return bal, errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.GetBalance.Scan(): ")
	}

	return bal, nil
}

func (e *pgMoneyManagerRepo) TransferFundsUserToUser(ctx context.Context,
	usrFrom entity.UserId, usrTo entity.UserId, fnd entity.Fund) error {

	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.TransferFundsUserToUser.Begin()")
	}

	//subtract amount from usrFrom
	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := DebitFundsFromUserSqlCmd

	stmtUser, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.TransferFundsUserToUser.Prepare() for DebitFundsFromUserSqlCmd")
	}

	if _, err := ts.Exec(ctx, stmtUser.SQL, usrFrom, fnd); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.TransferFundsUserToUser.Exec() for DebitFundsFromUserSqlCmd")
	}

	//add amount to usrTo
	prepPart = fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd = AddFundsToUserSqlCmd

	stmtReserve, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.TransferFundsUserToUser.Prepare() for AddFundsToReserveSqlCmd")
	}

	if _, err := ts.Exec(ctx, stmtReserve.SQL, usrTo, fnd); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.TransferFundsUserToUser.Exec() for AddFundsToReserveSqlCmd")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.TransferFundsUserToUser.Commit()")
	}

	return nil
}

func (e *pgMoneyManagerRepo) DebitFunds(ctx context.Context, usrId entity.UserId, fnd entity.Fund) error {

	ts, err := e.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.DebitFunds.Begin()")
	}

	prepPart := fmt.Sprintf("%d", time.Now().UnixNano())
	sqlCmd := DebitFundsFromUserSqlCmd

	stmt, err := ts.Prepare(ctx, "add"+prepPart, sqlCmd)
	if err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.DebitFunds.Prepare()")
	}

	if _, err := ts.Exec(ctx, stmt.SQL, usrId, fnd); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.DebitFunds.Exec()")
	}

	if err := ts.Commit(ctx); err != nil {
		return errors.Wrap(err, "MoneyManager.pgMoneyManagerRepo.DebitFunds.Commit()")
	}

	return nil
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
