package money_manager

import (
	"context"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

func (e *moneyManager) ReserveMoney(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error {
	if !isValidReserveIds(res) {
		return errors.New("err in moneyManager.ReserveMoney(): Invalid id in reserve")
	}

	fndToReserve, err := makeMoney(fndVal, fndUnit)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.ReserveMoney.makeMoney():")
	}

	usrBal, err := e.repo.GetBalance(ctx, res.UserId)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.ReserveMoney.GetBalance():")
	}
	if !isValidReserveOperation(usrBal, fndToReserve) {
		return errors.New("err in moneyManager.ReserveMoney(): Invalid reserve operation")
	}

	res.Amount = fndToReserve
	return e.repo.ReserveMoney(ctx, res)
}

func (e *moneyManager) RevokeReserve(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error {
	return nil
}

func (e *moneyManager) AcceptReserve(ctx context.Context, res entity.Reserve, funVal string, funUnit string) error {
	return nil
}

func isValidReserveIds(res entity.Reserve) bool {
	return isValidUserId(res.UserId) && res.OrderId != "" && res.ServiceId != ""
}

func isValidReserveOperation(bal entity.Balance, toReserve entity.Fund) bool {
	isHaveMoney := bal.Available >= toReserve
	curReserve := bal.Current - bal.Available

	return isHaveMoney && isValidFundSum(curReserve, toReserve)
}
