package money_manager

import (
	"context"
)

func (e *moneyManager) HealthCheck(ctx context.Context) error {
	return e.repo.HealthCheck(ctx)
}
