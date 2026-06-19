package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Character struct {
	ID            uuid.UUID       `db:"id"`
	UserID        uuid.UUID       `db:"user_id"`
	SystemName    string          `db:"system_name"`
	CharacterType string          `db:"character_type"`
	Name          string          `db:"name"`
	Data          json.RawMessage `db:"data"`
	CreatedAt     time.Time       `db:"created_at"`
}
