package repository

import (
	"time"

	"github.com/google/uuid"
)

type RollHistory struct {
	ID              uuid.UUID  `db:"id"`
	SystemName      string     `db:"system_name"`
	ActionType      string     `db:"action_type"`
	RequestPayload  any        `db:"request_payload"`
	ResponsePayload any        `db:"response_payload" `
	CampaignID      *uuid.UUID `db:"campaign_id"`
	CharacterID     *uuid.UUID `db:"character_id"`
	CreatedAt       time.Time  `db:"created_at"`
}
