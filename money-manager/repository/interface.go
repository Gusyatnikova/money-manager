package repository

import (
	"context"
	"money-manager/money-manager/entity"
)

type MoneyManagerRepo interface {
	AddUser(context.Context, entity.UserId, entity.MoneyAmount) error
	AddMoneyToUser(context.Context, entity.UserId, entity.MoneyAmount) error
	TransferMoneyUserToUser(context.Context, entity.UserId, entity.UserId, entity.MoneyAmount) error
	DebitMoney(context.Context, entity.UserId, entity.MoneyAmount) error
	GetBalance(context.Context, entity.UserId) (entity.Balance, error)

	ReserveMoney(context.Context, entity.Reserve) error
	GetReserve(context.Context, entity.Reserve) (entity.Reserve, error)
	CancelReserve(context.Context, entity.Reserve) error
	AcceptReserve(context.Context, entity.Reserve) error

	MakeReportMoneyPerService(ctx context.Context, reportPeriod entity.ReportPeriod) (entity.ReportMoneyPerService, error)

	HealthCheck(ctx context.Context) error
}
