package app

import (
	"context"
	"fmt"
	v1 "github.com/erminson/auth-var/internal/controller/http/v1"
	"github.com/erminson/auth-var/pkg/httpserver"
	"github.com/julienschmidt/httprouter"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	fmt.Println("App running...")

	ctx := context.Background()
	router := httprouter.New()
	v1.NewRouter(ctx, router /* logger and usecases */)
	httpServer := httpserver.New(router, httpserver.Port("8088"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-httpServer.Notify():
		fmt.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case s := <-interrupt:
		fmt.Println(fmt.Errorf("app - Run - signal: %s", s.String()))
	}

	err := httpServer.Shutdown()
	if err != nil {
		fmt.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
