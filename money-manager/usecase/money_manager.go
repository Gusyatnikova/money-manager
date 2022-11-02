package usecase

type moneyManager struct {
}

func NewMoneyManagerUseCase() MoneyManagerUseCase {
	return &moneyManager{}
}

func (e *moneyManager) AddFunds() {}
func (e *moneyManager) GetFunds() {}
