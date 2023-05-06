package user_handler

import (
	"context"
	"fmt"
	"net/http"
)

type UserHandler struct {
	// logger
	// user usecase
}

func New( /* logger and usecase */ ) *UserHandler {
	return &UserHandler{ /* logger and usecase */ }
}

func (h *UserHandler) CreateUser(ctx context.Context) http.Handler {
	fn := func(http.ResponseWriter, *http.Request) {
		fmt.Println("POST: /api/v1/user")
	}

	return http.HandlerFunc(fn)
}

func (h *UserHandler) GetUsers(ctx context.Context) http.Handler {
	fn := func(http.ResponseWriter, *http.Request) {
		fmt.Println("GET: /api/v1/user")
	}

	return http.HandlerFunc(fn)
}

func (h *UserHandler) GetUserById(ctx context.Context) http.Handler {
	fn := func(http.ResponseWriter, *http.Request) {
		fmt.Println("GET: /api/v1/user/:id")
	}

	return http.HandlerFunc(fn)
}

func (h *UserHandler) DeleteUserById(ctx context.Context) http.Handler {
	fn := func(http.ResponseWriter, *http.Request) {
		fmt.Println("DELETE: /api/v1/user/:id")
	}

	return http.HandlerFunc(fn)
}
