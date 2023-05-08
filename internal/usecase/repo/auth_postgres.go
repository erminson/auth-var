package repo

import (
	"context"
	"github.com/erminson/auth-var/internal/entity"
	"github.com/erminson/auth-var/pkg/postgres"
)

type AuthRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (r *AuthRepo) Store(ctx context.Context, ap entity.AuthPhone) error {
	sql := `
		INSERT INTO public.auth_phone (phone, code, created_at, valid_till)
		VALUES ($1, $2, $3, $4);
	`

	rows, err := r.Pool.Query(ctx, sql, ap.Phone, ap.Code, ap.CreatedAt, ap.ValidTil)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (r *AuthRepo) GetLastByPhoneNumber(ctx context.Context, phoneNumber string) (entity.AuthPhone, error) {
	sql := `
		SELECT id, phone, code, created_at, valid_till, confirmed_at
		FROM public.auth_phone
		WHERE phone LIKE $1
		ORDER BY created_at DESC
		LIMIT 1;
	`

	row := r.Pool.QueryRow(ctx, sql, phoneNumber)
	var ent entity.AuthPhone
	if err := row.Scan(&ent.Id, &ent.Phone, &ent.Code, &ent.CreatedAt, &ent.ValidTil, &ent.ConfirmedAt); err != nil {
		return entity.AuthPhone{}, err
	}

	return ent, nil
}

func (r *AuthRepo) ConfirmPhoneNumberById(ctx context.Context, id int) (int64, error) {
	sql := `
		UPDATE public.auth_phone
		SET confirmed_at = now()
		WHERE id = $1 AND confirmed_at IS NULL
	`

	res, err := r.Pool.Exec(ctx, sql, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}
