package usecase

import (
	"context"

	"money-manager/money-manager/entity"
)

type moneyManager struct {
	repo MoneyManagerRepo
}

func NewMoneyManagerUseCase(repo MoneyManagerRepo) MoneyManagerUseCase {
	return &moneyManager{
		repo: repo,
	}
}

func (e *moneyManager) AddFundsToUser(ctx context.Context, balanceOp entity.BalanceOperation) error {
	//check if userId is not empty and others
	return e.repo.AddFundsToUser(ctx, balanceOp)
}

func (e *moneyManager) GetBalance(ctx context.Context, usr entity.User) (entity.UserBalance, error) {
	return e.repo.GetBalance(ctx, usr)
}
