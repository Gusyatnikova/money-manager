package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"money-manager/money-manager/app/config"
)

func Connection(ctx context.Context, pgCfg config.PG) (*pgxpool.Pool, error) {
	confStr := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s port=%d pool_max_conns=%d",
		pgCfg.User, pgCfg.Password, pgCfg.Host, pgCfg.DbName, pgCfg.Port, pgCfg.PoolMax)

	conf, err := pgxpool.ParseConfig(confStr)
	if err != nil {
		return nil, err
	}

	pgConn, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	return pgConn, nil
}
