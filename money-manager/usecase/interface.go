package usecase

import (
	"context"

	"money-manager/money-manager/entity"
)

type MoneyManagerUseCase interface {
	AddMoneyToUser(ctx context.Context, usrId entity.UserId, money entity.Money) error
	TransferMoneyUserToUser(ctx context.Context, usrSrcId entity.UserId, usrDstId entity.UserId, money entity.Money) error
	DebitMoney(ctx context.Context, usrId entity.UserId, money entity.Money) error
	GetBalance(ctx context.Context, usrId entity.UserId) (entity.Balance, error)

	ReserveMoney(ctx context.Context, res entity.Reserve, money entity.Money) error
	AcceptReserve(ctx context.Context, res entity.Reserve, money entity.Money) error
	CancelReserve(ctx context.Context, res entity.Reserve) error

	HealthCheck(ctx context.Context) error
}
