package usecase

import (
	"context"
	"github.com/erminson/auth-var/internal/entity"
	"github.com/erminson/auth-var/internal/usecase/repo"
)

type Auth struct {
	repo *repo.AuthRepo
}

func New(r *repo.AuthRepo) *Auth {
	return &Auth{repo: r}
}

func (a *Auth) GenerateConfirmationCode(ctx context.Context, phoneNumber string) (entity.Message, error) {
	return entity.Message{}, nil
}

func (a *Auth) ConfirmPhoneNumber(ctx context.Context, phoneNumber, code string) (entity.Tokens, error) {
	return entity.Tokens{}, nil
}
