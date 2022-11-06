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
	"money-manager/money-manager/delivery"
	"money-manager/money-manager/repository/postgres"
	"money-manager/money-manager/usecase"
	"money-manager/money-manager/usecase/money_manager"
	"money-manager/pkg/http"
	dbConn "money-manager/pkg/repository/postgres"
)

type moneyManager struct {
	httpServer http.Server
	usecase    usecase.MoneyManagerUseCase
}

func NewMoneyManager(ctx context.Context) delivery.MoneyManagerServer {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic().Msgf("err in NewMoneyManager.NewConfig(): %s", err.Error())
	}

	initLogger()

	pgConn, err := dbConn.Connection(ctx, cfg.Pg)
	if err != nil {
		log.Panic().Msgf("err in NewMoneyManager.dbConn.Connection(): %s", err.Error())
	}

	moneyManagerRepo := postgres.NewPgMoneyManagerRepo(pgConn)
	uc := money_manager.NewMoneyManagerUseCase(moneyManagerRepo)

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
