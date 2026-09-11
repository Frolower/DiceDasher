package repository

import (
	"backend/internal/auth"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ pool *pgxpool.Pool }

var _ auth.Store = (*Repository)(nil)

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

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
