package v1

import (
	"context"
	auth_handler "github.com/erminson/auth-var/internal/controller/http/v1/auth"
	user_handler "github.com/erminson/auth-var/internal/controller/http/v1/user"
	"github.com/erminson/auth-var/internal/usecase"
	"github.com/erminson/auth-var/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(ctx context.Context, log logger.Interface, router *httprouter.Router, auth *usecase.Auth) {
	userHandler := user_handler.New( /* logger and usecase */ )
	router.POST("/api/v1/user", userHandler.CreateUser(ctx))
	router.GET("/api/v1/user", userHandler.GetUsers(ctx))
	router.GET("/api/v1/user/:id", userHandler.GetUserById(ctx))
	router.DELETE("/api/v1/user/:id", userHandler.DeleteUserById(ctx))

	authHandler := auth_handler.New(log, auth)
	router.POST("/api/v1/auth/phone", authHandler.GenerateConfirmationCode(ctx))
	router.POST("/api/v1/auth/confirm", authHandler.ConfirmPhoneNumber(ctx))
	router.POST("/api/v1/auth/refresh", authHandler.Refresh(ctx))
}
