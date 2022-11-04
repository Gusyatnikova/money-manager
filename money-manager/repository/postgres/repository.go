package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/oklog/ulid/v2"

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

func generateID() string {
	return ulid.Make().String()
}
