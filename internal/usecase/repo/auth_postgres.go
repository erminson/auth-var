package repo

import "github.com/erminson/auth-var/pkg/postgres"

type AuthRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}
