package money_manager

import (
	"context"
	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
	"money-manager/money-manager/repository"
	"money-manager/money-manager/usecase"
)

type moneyManager struct {
	repo repository.MoneyManagerRepo
}

func NewMoneyManagerUseCase(repo repository.MoneyManagerRepo) usecase.MoneyManagerUseCase {
	return &moneyManager{
		repo: repo,
	}
}

func (e *moneyManager) AddMoneyToUser(ctx context.Context, usrId entity.UserId, fndVal string, fndUnit string) error {
	if !isValidUserId(usrId) {
		return usecase.ErrInvalidUser
	}

	fndToAdd, err := makeMoney(fndVal, fndUnit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}

	//todo:use returns in database
	userBal, err := e.repo.GetBalance(ctx, usrId)
	if err != nil {
		if replaceNotFoundErr(err) != nil {
			err = e.repo.AddUser(ctx, usrId, fndToAdd)
		}
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.AddMoneyToUser().repo.GetBalance:")
	}

	if !isValidFundSum(balanceToFund(userBal), fndToAdd) {
		return usecase.ErrMoneyLimitIsExceeded
	}

	//todo: money, err := e.repo.AddMoneyToUser(ctx, usrId, fndToAdd)
	err = e.repo.AddMoneyToUser(ctx, usrId, fndToAdd)

	return errors.Wrap(err, "Err in AddMoneyToUser.repo.AddMoneyToUser:")
}

func (e *moneyManager) GetBalance(ctx context.Context, usrId entity.UserId) (entity.Balance, error) {
	bal := entity.Balance{}

	if !isValidUserId(usrId) {
		return bal, usecase.ErrInvalidUser
	}

	bal, err := e.repo.GetBalance(ctx, usrId)

	return bal, errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.GetBalance().repo.GetBalance:")
}

func (e *moneyManager) DebitMoney(ctx context.Context, usrId entity.UserId, fndVal string, fndUnit string) error {
	if !isValidUserId(usrId) {
		return usecase.ErrInvalidUser
	}

	fndToDebit, err := makeMoney(fndVal, fndUnit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}

	userBal, err := e.repo.GetBalance(ctx, usrId)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.DebitMoney.GetBalance():")
	}

	if !isValidDebit(userBal.Available, fndToDebit) {
		return usecase.ErrNotEnoughMoney
	}

	//todo: money, err := e.repo.DebitMoney(ctx, usrId, fndToDebit)
	err = e.repo.DebitMoney(ctx, usrId, fndToDebit)

	return errors.Wrap(err, "err in err in moneyManager.DebitMoney().repo.DebitMoney:")
}

func (e *moneyManager) TransferMoneyUserToUser(ctx context.Context,
	usrFromId entity.UserId, usrToId entity.UserId, fndVal string, fndUnit string) error {

	if usrFromId == usrToId || !isValidUserId(usrFromId) || !isValidUserId(usrToId) {
		return usecase.ErrInvalidUser
	}

	fndToTransfer, err := makeMoney(fndVal, fndUnit)
	if err != nil {
		return usecase.ErrInvalidMoney
	}

	srcBal, err := e.GetBalance(ctx, usrFromId)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.TransferMoneyUserToUser.GetBalance():")
	}

	if !isValidDebit(srcBal.Available, fndToTransfer) {
		return usecase.ErrNotEnoughMoney
	}

	dstBal, err := e.GetBalance(ctx, usrToId)
	if err != nil {
		return errors.Wrap(replaceNotFoundErr(err), "err in moneyManager.TransferMoneyUserToUser.GetBalance():")
	}

	if !isValidFundSum(balanceToFund(dstBal), fndToTransfer) {
		return usecase.ErrMoneyLimitIsExceeded
	}

	//todo: returns
	err = e.repo.TransferMoneyUserToUser(ctx, usrFromId, usrToId, fndToTransfer)

	return errors.Wrap(err, "err in err in moneyManager.TransferMoneyUserToUser().repo.TransferMoneyUserToUser:")
}

//Предполагается выполнять запрос к мастер-системе юзеров
func isValidUserId(usr entity.UserId) bool {
	return string(usr) != ""
}

//replaceNotFoundErr checks if err is repository.ErrNotFound then return usecase.ErrNotFound
//otherwise returns err
func replaceNotFoundErr(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return usecase.ErrNotFound
	}

	return err
}
