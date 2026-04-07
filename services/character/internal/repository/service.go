package repository

import (
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

func (r *Repository) InsertCharacter(ctx context.Context, ) (uuid.UUID, error) {
	const q = `
INSERT INTO characters ()
`

	var id uuid.UUID

	err := r.pool.QueryRow()

	return
}
