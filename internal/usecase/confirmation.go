package usecase

import (
	"context"
	"github.com/erminson/auth-var/internal/entity"
)

type Confirmation struct {
}

func New() *Confirmation {
	return &Confirmation{}
}

func (a *Confirmation) GenerateConfirmationCode(ctx context.Context, phoneNumber string) (entity.Message, error) {
	return entity.Message{}, nil
}

func (a *Confirmation) ConfirmPhoneNumber(ctx context.Context, phoneNumber, code string) (entity.Tokens, error) {
	return entity.Tokens{}, nil
}
