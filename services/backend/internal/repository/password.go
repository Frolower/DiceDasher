package repository

import (
	"backend/internal/user"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) UpdatePassword(ctx context.Context, id uuid.UUID, oldHash, newHash string) error {
	const q = `UPDATE public.users SET hash = $1 WHERE id = $2 AND hash = $3;`
	tag, err := r.pool.Exec(ctx, q, newHash, id, oldHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *Repository) GetHashByID(ctx context.Context, id uuid.UUID) (string, error) {
	var hash string

	const q = `SELECT hash FROM public.users WHERE id = $1;`
	err := r.pool.QueryRow(ctx, q, id).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", user.ErrUserNotFound
	}
	if err != nil {
		return "", err
	}

	return hash, nil
}
