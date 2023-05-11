package auth_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/erminson/auth-var/internal/usecase"
	"github.com/erminson/auth-var/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"regexp"
)

var (
	ErrPhoneNumber = errors.New("invalid phone number")
	ErrConfirmCode = errors.New("invalid confirm code")
)

type AuthPhoneRequest struct {
	Phone string `json:"phone"`
}

type PhoneConfirmRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type RefreshRequest struct {
	Refresh string `json:"refresh"`
}

type MessageResponse struct {
	Text   string `json:"text"`
	Digit  int    `json:"digit"`
	Source string `json:"source"`
}

type TokensResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type AccessResponse struct {
	Access string `json:"access"`
}

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type AuthHandler struct {
	l logger.Interface
	u *usecase.Auth
}

func New(log logger.Interface, uc *usecase.Auth) *AuthHandler {
	return &AuthHandler{
		l: log,
		u: uc,
	}
}

func (h *AuthHandler) GenerateConfirmationCode(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}
			h.errorResponse(w, apiError)
			return
		}

		var authPhone AuthPhoneRequest
		err = json.Unmarshal(data, &authPhone)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}
			h.errorResponse(w, apiError)
			return
		}

		err = h.validateAuthPhoneRequest(&authPhone)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{
				Status: http.StatusBadRequest,
				Title:  fmt.Sprintf("Bad Request. Error: %s", err.Error()),
			}
			h.errorResponse(w, apiError)
			return
		}

		msg, err := h.u.GenerateConfirmationCode(ctx, authPhone.Phone)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{
				Status: http.StatusBadRequest,
				Title:  fmt.Sprintf("Bad Request. Error: %s", err.Error()),
			}
			h.errorResponse(w, apiError)
			return
		}

		m := MessageResponse{
			Text:   msg.Text,
			Digit:  msg.Digit,
			Source: string(msg.Source),
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			h.l.Fatal(err.Error())
		}
	}

	return fn
}

func (h *AuthHandler) ConfirmPhoneNumber(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}
			h.errorResponse(w, apiError)
			return
		}

		var phoneConfirm PhoneConfirmRequest
		err = json.Unmarshal(data, &phoneConfirm)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}
			h.errorResponse(w, apiError)
			return
		}

		err = h.validatePhoneConfirmRequest(&phoneConfirm, 4)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{
				Status: http.StatusBadRequest,
				Title:  fmt.Sprintf("Bad Request. Error: %s", err.Error()),
			}
			h.errorResponse(w, apiError)
			return
		}

		t, err := h.u.ConfirmPhoneNumber(ctx, phoneConfirm.Phone, phoneConfirm.Code)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{
				Status: http.StatusBadRequest,
				Title:  fmt.Sprintf("Bad Request. Error: %s", err.Error()),
			}
			h.errorResponse(w, apiError)
			return
		}

		m := TokensResponse{
			Access:  t.Access,
			Refresh: t.Refresh,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			h.l.Fatal(err.Error())
		}
	}

	return fn
}

func (h *AuthHandler) Refresh(ctx context.Context) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}
			h.errorResponse(w, apiError)
			return
		}

		var refreshRequest RefreshRequest
		err = json.Unmarshal(data, &refreshRequest)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{Status: http.StatusBadRequest, Title: "Bad Request"}
			h.errorResponse(w, apiError)
			return
		}

		token, err := h.u.Refresh(ctx, refreshRequest.Refresh, 1)
		if err != nil {
			h.l.Error(err.Error())
			apiError := &ApiError{
				Status: http.StatusBadRequest,
				Title:  fmt.Sprintf("Bad Request. Error: %s", err.Error()),
			}
			h.errorResponse(w, apiError)
			return
		}

		a := AccessResponse{
			Access: token.Token,
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(a); err != nil {
			h.l.Fatal(err.Error())
		}
	}

	return fn
}

func (h *AuthHandler) errorResponse(w http.ResponseWriter, apiError *ApiError) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(apiError.Status)

	response := JsonErrorResponse{Error: apiError}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.l.Fatal(err.Error())
	}
}

func (h *AuthHandler) validateAuthPhoneRequest(in *AuthPhoneRequest) error {
	if len(in.Phone) == 0 {
		return ErrPhoneNumber
	}

	re := regexp.MustCompile(`^\d+$`)
	if !re.MatchString(in.Phone) {
		return ErrPhoneNumber
	}

	return nil
}

func (h *AuthHandler) validatePhoneConfirmRequest(in *PhoneConfirmRequest, d int) error {
	if len(in.Phone) == 0 {
		return ErrPhoneNumber
	}

	re := regexp.MustCompile(`^\d+$`)
	if !re.MatchString(in.Phone) {
		return ErrPhoneNumber
	}

	if len(in.Code) != d {
		return ErrConfirmCode
	}

	if !re.MatchString(in.Code) {
		return ErrConfirmCode
	}

	return nil
}
