package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"money-manager/money-manager/repository"
)

type pgMoneyManagerRepo struct {
	db *pgxpool.Pool
}

func NewPgMoneyManagerRepo(db *pgxpool.Pool) repository.MoneyManagerRepo {
	return &pgMoneyManagerRepo{
		db: db,
	}
}

// RunTx exec sql with transaction
func (e *pgMoneyManagerRepo) RunTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	var err error
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			// a panic occurred, rollback and repanic
			_ = tx.Rollback(ctx)
			panic(p)
		case err != nil:
			// something went wrong, rollback
			_ = tx.Rollback(ctx)
		default:
			// all good, commit
			err = tx.Commit(ctx)
		}
	}()
	err = fn(tx)
	return err
}
