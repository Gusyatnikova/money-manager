package delivery

type MoneyManagerServer interface {
	Run()
	Shutdown()
	ListenForShutdown()
}
