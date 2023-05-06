package v1

import (
	"context"
	user_handler "github.com/erminson/auth-var/internal/controller/http/v1/user"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(ctx context.Context, router *httprouter.Router /* logger and usecases */) {
	userHandler := user_handler.New( /* logger and usecase */ )
	router.HandlerFunc("POST", "/api/v1/user", userHandler.CreateUser(ctx).ServeHTTP)
	router.HandlerFunc("GET", "/api/v1/user", userHandler.GetUsers(ctx).ServeHTTP)
	router.HandlerFunc("GET", "/api/v1/user/:id", userHandler.GetUserById(ctx).ServeHTTP)
	router.HandlerFunc("DELETE", "/api/v1/user/:id", userHandler.DeleteUserById(ctx).ServeHTTP)
}
