package usecase

import (
	"context"

	"money-manager/money-manager/entity"
)

type MoneyManagerServer interface {
	Run()
	Shutdown()
	ListenForShutdown()
}

type MoneyManagerUseCase interface {
	AddFundsToUser(ctx context.Context, usr entity.User, funValue string, funUnit string) error
	GetBalance(ctx context.Context, usr entity.User) (entity.Balance, error)
}

type MoneyManagerRepo interface {
	AddUser(context.Context, entity.User, entity.Fund) error
	AddFundsToUser(context.Context, entity.User, entity.Fund) error
	GetBalance(context.Context, entity.User) (entity.Balance, error)
}
