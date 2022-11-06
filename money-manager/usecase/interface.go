package usecase

import (
	"context"

	"money-manager/money-manager/entity"
)

type MoneyManagerUseCase interface {
	AddMoneyToUser(ctx context.Context, usrId entity.UserId, fndValue string, fndUnit string) error
	TransferMoneyUserToUser(ctx context.Context, usrSrcId entity.UserId, usrDstId entity.UserId, fndVal string, fndUnit string) error
	DebitMoney(ctx context.Context, usrId entity.UserId, fndVal string, fndUnit string) error
	GetBalance(ctx context.Context, usrId entity.UserId) (entity.Balance, error)

	ReserveMoney(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error
	AcceptReserve(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error
	RevokeReserve(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error
}
