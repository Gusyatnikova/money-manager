package usecase

import (
	"context"

	"github.com/pkg/errors"

	"money-manager/money-manager/entity"
)

const (
	RUB string = "RUB"
	KOP        = "KOP"
)

const (
	KopVal int64 = 1
	RubVal       = 100 * KopVal
)

type moneyManager struct {
	repo MoneyManagerRepo
}

type Fund struct {
	Unit    string
	FundStr string
	fund    entity.Fund
}

func NewMoneyManagerUseCase(repo MoneyManagerRepo) MoneyManagerUseCase {
	return &moneyManager{
		repo: repo,
	}
}

func (e *moneyManager) AddFundsToUser(ctx context.Context, usr entity.User, fndVal string, fndUnit string) error {
	if !e.isValidUser(usr) {
		return errors.New("err in moneyManager.AddFundsToUser(): Invalid user")
	}

	fndToAdd, err := e.makeFunds(fndVal, fndUnit)
	if err != nil {
		return errors.Wrap(err, "err in moneyManager.AddFundsToUser.makeFunds():")
	}

	userBal, err := e.repo.GetBalance(ctx, usr)
	if err != nil {
		return e.repo.AddUser(ctx, usr, fndToAdd)
	}

	if !e.isValidFundSum(e.balanceToFund(userBal), fndToAdd) {
		return errors.New("err in moneyManager.AddFundsToUser(): result user's funds is too big")
	}

	return e.repo.AddFundsToUser(ctx, usr, fndToAdd)
}

func (e *moneyManager) GetBalance(ctx context.Context, usr entity.User) (entity.Balance, error) {
	bal := entity.Balance{}

	if !e.isValidUser(usr) {
		return bal, errors.New("err in moneyManager.GetBalance(): user is invalid")
	}

	return e.repo.GetBalance(ctx, usr)
}
