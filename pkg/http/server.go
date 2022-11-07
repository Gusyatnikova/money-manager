package http

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"money-manager/money-manager/delivery/http/http_v1"
	"money-manager/money-manager/usecase"
)

type Server interface {
	Run()
	Shutdown()
}

type ServerConfig struct {
	Address string
}

type server struct {
	ctx        context.Context
	httpServer *echo.Echo
}

func NewServer(ctx context.Context, cfg ServerConfig, uc usecase.MoneyManagerUseCase) Server {
	e := echo.New()
	e.Server.Addr = cfg.Address

	e.Use(
		mw.LoggerWithConfig(mw.LoggerConfig{
			Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}","method":"${method}",` +
				`"uri":"${uri}","query":"${query}","status":${status},"error":"${error}"}` + "\n",
		}),
		mw.Recover())

	http_v1.NewServerHandler(e, uc)

	return &server{
		ctx:        ctx,
		httpServer: e,
	}
}

func (e *server) Run() {
	log.Info().Msgf("HTTP server listening at %v", e.httpServer.Server.Addr)

	if err := e.httpServer.Server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panic().Msgf("err in server.Run(): %s", err.Error())
		}
	}
}

func (e *server) Shutdown() {
	ctx, cancel := context.WithTimeout(e.ctx, 2*time.Second)
	defer cancel()

	log.Info().Msgf("Shutting down the server")

	if err := e.httpServer.Server.Shutdown(ctx); err != nil {
		log.Error().Msgf("err in server.Shutdown(): %s", err.Error())
	}
}
