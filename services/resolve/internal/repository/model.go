package repository

import (
	"encoding/json/jsontext"
	"time"

	"github.com/google/uuid"
)

type RollHistory struct {
	ID              uuid.UUID      `db:"id"`
	RequestID       uuid.UUID      `db:"request_id"`
	SystemName      string         `db:"system_name"`
	ActionType      string         `db:"action_type"`
	StatePayload    jsontext.Value `db:"state_payload"`
	RequestPayload  jsontext.Value `db:"request_payload"`
	ResponsePayload jsontext.Value `db:"response_payload" `
	CampaignID      *uuid.UUID     `db:"campaign_id"`
	CharacterID     *uuid.UUID     `db:"character_id"`
	CreatedAt       time.Time      `db:"created_at"`
}
