package repo

import (
	"context"
	"errors"
	"github.com/erminson/auth-var/internal/entity"
	"github.com/erminson/auth-var/pkg/postgres"
)

type AuthRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (r *AuthRepo) StoreAuth(ctx context.Context, ap entity.AuthPhone) error {
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

func (r *AuthRepo) GetLastAuthByPhoneNumber(ctx context.Context, phoneNumber string) (entity.AuthPhone, error) {
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

func (r *AuthRepo) ConfirmPhoneNumberAndSaveUser(ctx context.Context, ap entity.AuthPhone) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sqlUpdateAuth := `
		UPDATE public.auth_phone
		SET confirmed_at = now()
		WHERE id = $1 AND confirmed_at IS NULL
	`

	res, err := tx.Exec(ctx, sqlUpdateAuth, ap.Id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("phone number not found")
	}

	sqlExistsUser := `
		SELECT EXISTS(
			SELECT id FROM public.user WHERE phone LIKE $1
		); 
	`
	sqlInsertUser := `
		INSERT INTO public.user (phone)
		VALUES ($1)
	`

	var exists bool
	row := tx.QueryRow(ctx, sqlExistsUser, ap.Phone)
	if err := row.Scan(&exists); err != nil {
		return err
	}

	if !exists {
		_, err := tx.Exec(ctx, sqlInsertUser, ap.Phone)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
