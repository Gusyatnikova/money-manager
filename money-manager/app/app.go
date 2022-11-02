package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/app/config"
	"money-manager/money-manager/usecase"
	"money-manager/pkg/http"
)

type moneyManager struct {
	httpServer http.Server
	usecase    usecase.MoneyManagerUseCase
}

func NewMoneyManager(ctx context.Context) usecase.MoneyManagerServer {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic().Msgf("err in NewMoneyManager.NewConfig(): %s", err.Error())
	}

	initLogger()

	uc := usecase.NewMoneyManagerUseCase()

	httpServerCfg := http.ServerConfig{
		Address: fmt.Sprint(cfg.Http.Host, ":", cfg.Http.Port),
	}
	httpServer := http.NewServer(ctx, httpServerCfg, uc)

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
