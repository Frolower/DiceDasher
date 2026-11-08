package repository

import (
	"encoding/json/jsontext"
	"time"

	"github.com/google/uuid"
)

type Character struct {
	ID            uuid.UUID      `db:"id"`
	UserID        uuid.UUID      `db:"user_id"`
	SystemName    string         `db:"system_name"`
	CharacterType string         `db:"character_type"`
	Name          string         `db:"name"`
	Data          jsontext.Value `db:"data"`
	CreatedAt     time.Time      `db:"created_at"`
}
