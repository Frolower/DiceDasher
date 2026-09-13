package repository

import (
	"backend/internal/auth"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateUser(ctx context.Context, rec auth.CreateRecord) (uuid.UUID, error) {
	const q = `INSERT INTO public.users (username, email, hash)
               VALUES ($1, $2, $3) RETURNING id`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, rec.Username, rec.Email, rec.PasswordHash).Scan(&id)
	if err != nil {
		return uuid.Nil, mapCreateError(err)
	}
	return id, nil
}

func mapCreateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && (pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_email_unique") {
			return auth.ErrAlreadyExists
		}
		return fmt.Errorf("insert user: PostgreSQL error %s", pgErr.Code)
	}
	return fmt.Errorf("insert user: %w", err)
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (auth.Credentials, error) {
	const q = `SELECT id, hash FROM public.users WHERE lower(username) = lower($1)`
	var credentials auth.Credentials
	err := r.pool.QueryRow(ctx, q, username).Scan(&credentials.UserID, &credentials.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.Credentials{}, auth.ErrUserNotFound
	}
	if err != nil {
		return auth.Credentials{}, fmt.Errorf("find user: %w", err)
	}
	return credentials, nil
}
