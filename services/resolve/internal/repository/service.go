package repository

import (
	"context"
	"errors"

	"diceDasher/pkg/dbutil"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")

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

func (r *Repository) InsertRollHistory(ctx context.Context, rec RollHistory) (uuid.UUID, error) {
	const q = `
INSERT INTO public.roll_history
(request_id, system_name, action_type, request_payload, response_payload, campaign_id, character_id, state_payload)
VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6, $7, $8::jsonb)
RETURNING id;
`
	var id uuid.UUID

	err := r.pool.QueryRow(
		ctx,
		q,
		rec.RequestID,
		rec.SystemName,
		rec.ActionType,
		rec.RequestPayload,
		rec.ResponsePayload,
		rec.CampaignID,
		rec.CharacterID,
		rec.StatePayload,
	).Scan(&id)

	return id, err
}

func (r *Repository) GetRollHistoryByID(ctx context.Context, id uuid.UUID) (RollHistory, error) {
	const q = `
SELECT id, system_name, action_type, request_payload, response_payload, state_payload
FROM public.roll_history
WHERE id = $1;
`
	var rec RollHistory
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&rec.ID,
		&rec.SystemName,
		&rec.ActionType,
		&rec.RequestPayload,
		&rec.ResponsePayload,
		&rec.StatePayload,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RollHistory{}, ErrNotFound
		}
		return RollHistory{}, err
	}
	return rec, nil
}
