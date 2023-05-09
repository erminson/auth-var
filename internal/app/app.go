package app

import (
	"context"
	"fmt"
	v1 "github.com/erminson/auth-var/internal/controller/http/v1"
	"github.com/erminson/auth-var/internal/usecase"
	"github.com/erminson/auth-var/internal/usecase/repo"
	"github.com/erminson/auth-var/internal/usecase/webapi"
	"github.com/erminson/auth-var/pkg/httpserver"
	jwt_client "github.com/erminson/auth-var/pkg/jwt"
	"github.com/erminson/auth-var/pkg/postgres"
	"github.com/julienschmidt/httprouter"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	fmt.Println("App running...")

	pg, err := postgres.New("postgresql://localhost:5430/authvar_db?user=authvar_dev_user&password=authvar_dev_paSSword")
	if err != nil {
		fmt.Println(fmt.Errorf("app - Run - postgres.New: %w", err).Error())
		return
	}
	defer pg.Close()

	authUseCase := usecase.New(
		repo.New(pg),
		webapi.New("token"),
		jwt_client.New("CIcaqLR27InWdldWaM96gXkPW90dc4tR8At3H7Sx"),
	)

	router := httprouter.New()
	v1.NewRouter(context.Background(), router, authUseCase)
	httpServer := httpserver.New(router, httpserver.Port("8088"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-httpServer.Notify():
		fmt.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case s := <-interrupt:
		fmt.Println(fmt.Errorf("app - Run - signal: %s", s.String()))
	}

	err = httpServer.Shutdown()
	if err != nil {
		fmt.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
