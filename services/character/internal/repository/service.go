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

func (r *Repository) InsertPlayerCreatedCharacter(ctx context.Context, rec Character) (uuid.UUID, error) {
	const q = `
INSERT INTO public.characters
(user_id, system_name, character_type, name, data)
VALUES ($1, $2, $3, $4, $5::jsonb)
RETURNING id;
`
	var id uuid.UUID

	err := r.pool.QueryRow(
		ctx,
		q,
		rec.UserID,
		rec.SystemName,
		rec.CharacterType,
		rec.Name,
		rec.Data,
	).Scan(&id)

	return id, err
}
