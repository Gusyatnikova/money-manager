package money_manager

import (
	"context"
	"money-manager/money-manager/usecase"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

func (e *moneyManager) ReserveMoney(ctx context.Context, resToAdd entity.Reserve, moneyVal string, moneyUnit string) error {
	if !isValidReserveIds(resToAdd) {
		return usecase.ErrInvalidReserve
	}

	usrBal, err := e.repo.GetBalance(ctx, resToAdd.UserId)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.ReserveMoney.GetBalance():")
	}

	moneyToReserve, err := makeMoney(moneyVal, moneyUnit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}
	if !isValidReserveOperation(usrBal, moneyToReserve) {
		return errors.New("err in moneyManager.ReserveMoney(): Invalid reserve operation")
	}

	resToAdd.Amount = moneyToReserve
	return e.repo.ReserveMoney(ctx, resToAdd)
}

func (e *moneyManager) CancelReserve(ctx context.Context, resToCancel entity.Reserve) error {
	if !isValidReserveIds(resToCancel) {
		return usecase.ErrInvalidReserve
	}

	_, err := e.repo.GetReserve(ctx, resToCancel)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.ReserveMoney.CancelReserve():")
	}

	err = e.repo.CancelReserve(ctx, resToCancel)

	return errors.Wrap(err, "err in err in moneyManager.CancelReserve().repo.CancelReserve:")
}

func (e *moneyManager) AcceptReserve(ctx context.Context, resToAccept entity.Reserve, moneyVal string, moneyUnit string) error {
	if !isValidReserveIds(resToAccept) {
		return usecase.ErrInvalidReserve
	}

	curReserve, err := e.repo.GetReserve(ctx, resToAccept)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.ReserveMoney.GetReserve():")
	}

	moneyToAccept, err := makeMoney(moneyVal, moneyUnit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}
	resToAccept.Amount = moneyToAccept

	if !isValidDebit(curReserve.Amount, moneyToAccept) {
		return usecase.ErrNotEnoughMoney
	}

	err = e.repo.AcceptReserve(ctx, resToAccept)

	return errors.Wrap(err, "err in err in moneyManager.AcceptReserve().repo.AcceptReserve:")
}

func isValidReserveIds(res entity.Reserve) bool {
	return isValidUserId(res.UserId) && res.OrderId != "" && res.ServiceId != ""
}

func isValidReserveOperation(curBal entity.Balance, toReserve entity.Fund) bool {
	isHaveMoney := curBal.Available >= toReserve
	curReserve := curBal.Current - curBal.Available

	return isHaveMoney && isValidFundSum(curReserve, toReserve)
}
