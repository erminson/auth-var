package v1

import (
	"context"
	auth_handler "github.com/erminson/auth-var/internal/controller/http/v1/auth"
	user_handler "github.com/erminson/auth-var/internal/controller/http/v1/user"
	"github.com/erminson/auth-var/internal/usecase"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(ctx context.Context, router *httprouter.Router /* logger and usecases */) {
	userHandler := user_handler.New( /* logger and usecase */ )
	router.POST("/api/v1/user", userHandler.CreateUser(ctx))
	router.GET("/api/v1/user", userHandler.GetUsers(ctx))
	router.GET("/api/v1/user/:id", userHandler.GetUserById(ctx))
	router.DELETE("/api/v1/user/:id", userHandler.DeleteUserById(ctx))

	confirmation := usecase.New()
	authHandler := auth_handler.New(confirmation)
	router.POST("/api/v1/auth/phone", authHandler.GenerateConfirmationCode(ctx))
	router.POST("/api/v1/auth/confirm", authHandler.ConfirmPhoneNumber(ctx))
}
