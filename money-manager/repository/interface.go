package repository

import (
	"context"

	"money-manager/money-manager/entity"
)

type MoneyManagerRepo interface {
	AddUser(context.Context, entity.UserId, entity.Fund) error
	AddMoneyToUser(context.Context, entity.UserId, entity.Fund) error
	TransferMoneyUserToUser(context.Context, entity.UserId, entity.UserId, entity.Fund) error
	DebitMoney(context.Context, entity.UserId, entity.Fund) error
	GetBalance(context.Context, entity.UserId) (entity.Balance, error)

	ReserveMoney(context.Context, entity.Reserve) error
	GetReserve(context.Context, entity.Reserve) (entity.Reserve, error)
	CancelReserve(context.Context, entity.Reserve) error
	AcceptReserve(context.Context, entity.Reserve) error

	HealthCheck(ctx context.Context) error
}
