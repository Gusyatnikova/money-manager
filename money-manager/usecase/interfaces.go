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
	AddFunds(context.Context, entity.BalanceOperation) error
	GetBalance(context.Context, entity.User) (entity.UserBalance, error)
}

type MoneyManagerRepoHandler interface {
	AddFunds(context.Context, entity.BalanceOperation) error
	GetBalance(context.Context, entity.User) (entity.UserBalance, error)
}
