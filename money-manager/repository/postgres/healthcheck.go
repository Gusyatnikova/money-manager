package postgres

import "context"

func (e *pgMoneyManagerRepo) HealthCheck(ctx context.Context) error {
	return e.db.Ping(ctx)
}
