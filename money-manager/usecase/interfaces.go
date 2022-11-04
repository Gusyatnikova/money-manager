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
	AddFundsToUser(ctx context.Context, usrId entity.UserId, fndValue string, fndUnit string) error
	GetBalance(ctx context.Context, usrId entity.UserId) (entity.Balance, error)

	Reserve(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error
	AcceptReserve(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error
}

type MoneyManagerRepo interface {
	AddUser(context.Context, entity.UserId, entity.Fund) error
	AddFundsToUser(context.Context, entity.UserId, entity.Fund) error
	GetBalance(context.Context, entity.UserId) (entity.Balance, error)

	AddReserve(ctx context.Context, res entity.Reserve) error
	AcceptReserve(ctx context.Context, resId entity.ReserveId) error
}
