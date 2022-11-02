package app

import (
	"context"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"

	"money-manager/money-manager/usecase"
	"money-manager/pkg/http"
)

type moneyManager struct {
	httpServer http.Server
	usecase    usecase.MoneyManagerUseCase
}

func NewMoneyManager(ctx context.Context) usecase.MoneyManagerServer {
	initLogger()

	uc := usecase.NewMoneyManagerUseCase()

	httpServer := http.NewServer(ctx)

	return &moneyManager{
		httpServer: httpServer,
		usecase:    uc,
	}
}

func (e *moneyManager) Run() {
	go e.httpServer.Run()
}

func (e *moneyManager) Shutdown() {
	e.httpServer.Shutdown()
}

func (e *moneyManager) ListenForShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	for v := range ch {
		if v == os.Interrupt {
			e.Shutdown()

			break
		}
	}
	close(ch)
}

func initLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
