package user_handler

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UserHandler struct {
	// logger
	// user usecase
}

func New( /* logger and usecase */ ) *UserHandler {
	return &UserHandler{ /* logger and usecase */ }
}

func (h *UserHandler) CreateUser(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("POST: /api/v1/user")
	}

	return fn
}

func (h *UserHandler) GetUsers(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("GET: /api/v1/user")
	}

	return fn
}

func (h *UserHandler) GetUserById(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Printf("GET: /api/v1/user/%s\n", ps.ByName("id"))
	}

	return fn
}

func (h *UserHandler) DeleteUserById(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("DELETE: /api/v1/user/:id")
	}

	return fn
}
