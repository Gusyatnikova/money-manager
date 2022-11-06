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

	ReserveMoney(ctx context.Context, res entity.Reserve) error
	RevokeReserve(ctx context.Context, res entity.Reserve) error
	AcceptReserve(ctx context.Context, resId entity.ReserveId) error
}
