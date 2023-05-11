package app

import (
	"context"
	"fmt"
	"github.com/erminson/auth-var/config"
	v1 "github.com/erminson/auth-var/internal/controller/http/v1"
	"github.com/erminson/auth-var/internal/usecase"
	"github.com/erminson/auth-var/internal/usecase/repo"
	"github.com/erminson/auth-var/internal/usecase/webapi"
	"github.com/erminson/auth-var/pkg/httpserver"
	jwt_client "github.com/erminson/auth-var/pkg/jwt"
	"github.com/erminson/auth-var/pkg/logger"
	"github.com/erminson/auth-var/pkg/postgres"
	"github.com/julienschmidt/httprouter"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New("debug")
	l.Info("App running...")

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	uc := usecase.New(
		repo.New(pg),
		webapi.New("token"),
		jwt_client.New(cfg.JWT.SecretKey),
	)

	// HTTP Server
	r := httprouter.New()
	v1.NewRouter(context.Background(), l, r, uc)
	httpServer := httpserver.New(r, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case s := <-interrupt:
		l.Error(fmt.Errorf("app - Run - signal: %s", s.String()))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
