package money_manager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
	"money-manager/money-manager/usecase"
)

type moneyManager struct {
	repo usecase.MoneyManagerRepo
}

func NewMoneyManagerUseCase(repo usecase.MoneyManagerRepo) usecase.MoneyManagerUseCase {
	return &moneyManager{
		repo: repo,
	}
}

func (e *moneyManager) AddFundsToUser(ctx context.Context, usr entity.UserId, fndVal string, fndUnit string) error {
	if !isValidUserId(usr) {
		return errors.New("err in moneyManager.AddFundsToUser(): Invalid user")
	}

	fndToAdd, err := makeFunds(fndVal, fndUnit)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.AddFundsToUser.makeFunds():")
	}

	userBal, err := e.repo.GetBalance(ctx, usr)
	if err != nil {
		return e.repo.AddUser(ctx, usr, fndToAdd)
	}

	if !isValidFundSum(balanceToFund(userBal), fndToAdd) {
		return errors.New("err in moneyManager.AddFundsToUser(): result user's funds is too big")
	}

	return e.repo.AddFundsToUser(ctx, usr, fndToAdd)
}

func (e *moneyManager) GetBalance(ctx context.Context, usr entity.UserId) (entity.Balance, error) {
	bal := entity.Balance{}

	if !isValidUserId(usr) {
		return bal, errors.New("err in moneyManager.GetBalance(): user is invalid")
	}

	return e.repo.GetBalance(ctx, usr)
}

func (e *moneyManager) Reserve(ctx context.Context, res entity.Reserve, fndVal string, fndUnit string) error {
	if !isValidReserveIds(res) {
		return errors.New("err in moneyManager.Reserve(): Invalid id in reserve")
	}

	fndToReserve, err := makeFunds(fndVal, fndUnit)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.Reserve.makeFunds():")
	}

	usrBal, err := e.repo.GetBalance(ctx, res.UserId)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.Reserve.GetBalance():")
	}
	if !isValidReserveOperation(usrBal, fndToReserve) {
		return errors.New("err in moneyManager.Reserve(): Invalid reserve operation")
	}

	res.Amount = fndToReserve
	return e.repo.AddReserve(ctx, res)
}

func (e *moneyManager) AcceptReserve(ctx context.Context, res entity.Reserve, funVal string, funUnit string) error {
	return nil
}

func (e *moneyManager) TransferFundsUserToUser(ctx context.Context,
	usrFromId entity.UserId, usrToId entity.UserId, fndVal string, fndUnit string) error {

	if !isValidUserId(usrFromId) || !isValidUserId(usrToId) {
		return errors.New("err in moneyManager.TransferFundsUserToUser(): Invalid user")
	}

	fndToTransfer, err := makeFunds(fndVal, fndUnit)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.DebitFunds.makeFunds():")
	}

	srcBal, err := e.GetBalance(ctx, usrFromId)
	if err != nil {
		return errors.New("err in moneyManager.TransferFundsUserToUser(): Invalid usrSrc balance")
	}

	if !isValidFundDebit(srcBal.Available, fndToTransfer) {
		return errors.New("err in moneyManager.TransferFundsUserToUser(): insufficient funds to pay")
	}

	_, err = e.GetBalance(ctx, usrToId)
	if err != nil {
		return errors.New(fmt.Sprintf("err in moneyManager.TransferFundsUserToUser(): no user with id: %s", usrToId))
	}

	return e.repo.TransferFundsUserToUser(ctx, usrFromId, usrToId, fndToTransfer)
}
func (e *moneyManager) DebitFunds(ctx context.Context, usrId entity.UserId, fndVal string, fndUnit string) error {

	if !isValidUserId(usrId) {
		return errors.New("err in moneyManager.DebitFunds(): Invalid user")
	}

	fndToDebit, err := makeFunds(fndVal, fndUnit)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.DebitFunds.makeFunds():")
	}

	userBal, err := e.repo.GetBalance(ctx, usrId)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.DebitFunds.GetBalance():")
	}

	if !isValidFundDebit(userBal.Available, fndToDebit) {
		return errors.New("err in moneyManager.DebitFunds(): insufficient funds")
	}

	return e.repo.DebitFunds(ctx, usrId, fndToDebit)
}
