package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/erminson/auth-var/internal/entity"
	"github.com/erminson/auth-var/internal/usecase/repo"
	"github.com/erminson/auth-var/internal/usecase/webapi"
	jwt_client "github.com/erminson/auth-var/pkg/jwt"
	"strings"
	"time"
)

var (
	ErrPhoneNumberNotFound        = errors.New("phone number not found")
	ErrConfirmationCodeIsWrong    = errors.New("confirmation code is wrong")
	ErrConfirmationCodeHasExpired = errors.New("confirmation code has expired")
)

type Auth struct {
	repo   *repo.AuthRepo
	webAPI *webapi.ConfirmAPI
	jwt    *jwt_client.JWTClient
}

func New(r *repo.AuthRepo, w *webapi.ConfirmAPI, j *jwt_client.JWTClient) *Auth {
	return &Auth{
		repo:   r,
		webAPI: w,
		jwt:    j,
	}
}

func (a *Auth) GenerateConfirmationCode(ctx context.Context, phoneNumber string) (entity.Message, error) {
	code, err := a.webAPI.GenerateCode(phoneNumber)
	if err != nil {
		return entity.Message{}, err
	}

	fmt.Println(code)

	m := entity.Message{
		Text:   fmt.Sprintf("Enter %d characters from %s", len(code), strings.ToUpper(entity.Sms.String())),
		Digit:  len(code),
		Source: entity.Sms,
	}

	createdAt := time.Now()
	authPhone := entity.AuthPhone{
		Phone:     phoneNumber,
		Code:      code,
		CreatedAt: createdAt,
		ValidTil:  createdAt.Add(time.Duration(time.Minute * 3)),
	}

	err = a.repo.Store(ctx, authPhone)
	if err != nil {
		return entity.Message{}, err
	}

	return m, nil
}

func (a *Auth) ConfirmPhoneNumber(ctx context.Context, phoneNumber, code string) (entity.Tokens, error) {
	authPhone, err := a.repo.GetLastByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return entity.Tokens{}, err
	}

	if authPhone.ConfirmedAt != nil {
		return entity.Tokens{}, ErrPhoneNumberNotFound
	}

	if authPhone.Code != code {
		return entity.Tokens{}, ErrConfirmationCodeIsWrong
	}

	if time.Now().After(authPhone.ValidTil) {
		return entity.Tokens{}, ErrConfirmationCodeHasExpired
	}

	count, err := a.repo.ConfirmPhoneNumberById(ctx, authPhone.Id)
	if err != nil {
		return entity.Tokens{}, err
	}

	if count == 0 {
		return entity.Tokens{}, ErrPhoneNumberNotFound
	}

	//err = a.webAPI.ConfirmNumber(phoneNumber, code)
	//if err != nil {
	//	return entity.Tokens{}, err
	//}

	td, err := a.jwt.GenerateTokenDetails(1)
	if err != nil {
		return entity.Tokens{}, err
	}

	return entity.Tokens{
		Access:  td.Access,
		Refresh: td.Refresh,
	}, nil
}
