package auth_handler

import (
	"context"
	"encoding/json"
	"github.com/erminson/auth-var/internal/entity"
	"github.com/erminson/auth-var/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

type AuthPhoneRequest struct {
	Phone string `json:"phone"`
}

type PhoneConfirmRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type MessageResponse struct {
	Text   string `json:"text"`
	Digit  byte   `json:"digit"`
	Source string `json:"source"`
}

type TokensResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int16  `json:"status"`
	Title  string `json:"title"`
}

type AuthHandler struct {
	u *usecase.Confirmation
}

func New(confirmation *usecase.Confirmation) *AuthHandler {
	return &AuthHandler{
		u: confirmation,
	}
}

func (h *AuthHandler) GenerateConfirmationCode(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}}
			if err = json.NewEncoder(w).Encode(response); err != nil {
				panic(err)
			}
			return
		}

		var authPhone AuthPhoneRequest
		err = json.Unmarshal(data, &authPhone)
		if err != nil {
			// logging

			w.WriteHeader(http.StatusBadRequest)
			response := JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}}
			if err = json.NewEncoder(w).Encode(response); err != nil {
				panic(err)
			}
			return
		}

		m := MessageResponse{
			Text:   "SMS",
			Digit:  4,
			Source: string(entity.Sms),
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(m); err != nil {
			panic(err)
		}
	}

	return fn
}

func (h *AuthHandler) ConfirmPhoneNumber(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}}
			if err = json.NewEncoder(w).Encode(response); err != nil {
				panic(err)
			}
			return
		}

		var phoneConfirm PhoneConfirmRequest
		err = json.Unmarshal(data, &phoneConfirm)
		if err != nil {
			// logging

			w.WriteHeader(http.StatusBadRequest)
			response := JsonErrorResponse{Error: &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}}
			if err = json.NewEncoder(w).Encode(response); err != nil {
				panic(err)
			}
			return
		}

		m := TokensResponse{
			Access:  "access_token",
			Refresh: "refresh_token",
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(m); err != nil {
			panic(err)
		}
	}

	return fn
}
