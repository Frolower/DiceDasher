package repository

import (
	"backend/internal/user"
	"context"
	"diceDasher/pkg/dbutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func FromContext(ctx context.Context) (*Repository, error) {
	base, err := dbutil.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return New(base.Pool()), nil
}

func (r *Repository) CreateUser(ctx context.Context, rec user.CreateRecord) (uuid.UUID, error) {
	const q = `
INSERT INTO public.users
(username, email, hash)
VALUES ($1, $2, $3)
RETURNING id;
`
	var id uuid.UUID

	err := r.pool.QueryRow(
		ctx,
		q,
		rec.Username,
		rec.Email,
		rec.PasswordHash,
	).Scan(&id)

	return id, err
}
