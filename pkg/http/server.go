package http

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Server interface {
	Run()
	Shutdown()
}

type server struct {
	ctx        context.Context
	httpServer *echo.Echo
}

func NewServer(ctx context.Context) Server {
	e := echo.New()

	e.Use(
		mw.LoggerWithConfig(mw.LoggerConfig{
			Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}","method":"${method}",` +
				`"uri":"${uri}","query":"${query}","status":${status},"error":"${error}"}` + "\n",
		}),
		mw.Recover())

	return &server{
		ctx:        ctx,
		httpServer: e,
	}
}

func (e *server) Run() {
	if err := e.httpServer.Server.ListenAndServe(); err != nil {
		log.Panic().Msgf("err in server.Run(): %s", err.Error())
	}

	log.Info().Msgf("HTTP server listening at %v", 1111)
}

func (e *server) Shutdown() {
	ctx, cancel := context.WithTimeout(e.ctx, 2*time.Second)
	defer cancel()

	log.Warn().Msgf("Shutting down the server")

	if err := e.httpServer.Server.Shutdown(ctx); err != nil {
		log.Error().Msgf("err in server.Shutdown(): %s", err.Error())
	}
}
