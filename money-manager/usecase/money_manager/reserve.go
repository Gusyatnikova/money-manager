package money_manager

import (
	"context"
	"money-manager/money-manager/usecase"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

func (e *moneyManager) ReserveMoney(ctx context.Context, resToAdd entity.Reserve, money entity.Money) error {
	if !isValidReserveIds(resToAdd) {
		return usecase.ErrInvalidReserve
	}

	usrBal, err := e.repo.GetBalance(ctx, resToAdd.UserId)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.ReserveMoney.GetBalance():")
	}

	moneyToReserve, err := makeMoneyAmount(money.Value, money.Unit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}

	if !isValidReserveOperation(usrBal, moneyToReserve) {
		return errors.New("err in moneyManager.ReserveMoney(): Invalid reserve operation")
	}

	_, err = e.repo.GetReserve(ctx, resToAdd)
	if err == nil {
		//reserve is already exists
		return errors.Wrap(usecase.ErrDuplication, "err in moneyManager.ReserveMoney:")
	}

	resToAdd.MoneyAmount = moneyToReserve
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

func (e *moneyManager) AcceptReserve(ctx context.Context, resToAccept entity.Reserve, money entity.Money) error {
	if !isValidReserveIds(resToAccept) {
		return usecase.ErrInvalidReserve
	}

	curReserve, err := e.repo.GetReserve(ctx, resToAccept)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.ReserveMoney.GetReserve():")
	}

	moneyToAccept, err := makeMoneyAmount(money.Value, money.Unit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}
	resToAccept.MoneyAmount = moneyToAccept

	if !isValidDebit(curReserve.MoneyAmount, moneyToAccept) {
		return usecase.ErrNotEnoughMoney
	}

	err = e.repo.AcceptReserve(ctx, resToAccept)

	return errors.Wrap(err, "err in err in moneyManager.AcceptReserve().repo.AcceptReserve:")
}

func isValidReserveIds(res entity.Reserve) bool {
	return isValidUserId(res.UserId) && res.OrderId != "" && res.ServiceId != ""
}

//isValidReserveOperation checks if curBal is more than toReserve
//and checks possible overflow after reservation
func isValidReserveOperation(curBal entity.Balance, toReserve entity.MoneyAmount) bool {
	isHaveMoney := curBal.Available >= toReserve
	usrCurReserve := curBal.Current - curBal.Available

	return isHaveMoney && isValidAmountSum(usrCurReserve, toReserve)
}
