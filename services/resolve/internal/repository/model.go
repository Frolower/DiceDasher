package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RollHistory struct {
	ID              uuid.UUID       `db:"id"`
	RequestID       uuid.UUID       `db:"request_id"`
	SystemName      string          `db:"system_name"`
	ActionType      string          `db:"action_type"`
	RequestPayload  json.RawMessage `db:"request_payload"`
	ResponsePayload json.RawMessage `db:"response_payload" `
	CampaignID      *uuid.UUID      `db:"campaign_id"`
	CharacterID     *uuid.UUID      `db:"character_id"`
	CreatedAt       time.Time       `db:"created_at"`
}
