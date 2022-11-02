package usecase

type MoneyManagerServer interface {
	Run()
	Shutdown()
	ListenForShutdown()
}

type MoneyManagerUseCase interface {
	AddFunds()
	GetBalance()
}
