package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (r *Repository) InsertRollHistory(ctx context.Context, rec RollHistory) (uuid.UUID, error) {
	const q = `
INSERT INTO public.roll_history
(request_id, system_name, action_type, request_payload, response_payload, campaign_id, character_id)
VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6, $7)
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
	).Scan(&id)

	return id, err
}

func (r *Repository) GetRollHistoryByID(ctx context.Context, id uuid.UUID) (RollHistory, error) {
	const q = `
SELECT request_payload, response_payload
FROM public.roll_history
WHERE id = $1;
`
	var rec RollHistory
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&rec.RequestPayload,
		&rec.ResponsePayload,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return RollHistory{}, errors.New("record not found")
		}
		return RollHistory{}, err
	}
	return rec, nil
}
