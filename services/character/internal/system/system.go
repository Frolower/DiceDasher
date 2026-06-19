package system

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type CreatedCharacter struct {
	UserID        uuid.UUID
	CharacterType string
	Name          string
	Data          json.RawMessage
}

type Character interface {
	CreateCharacter(ctx context.Context, raw json.RawMessage) (CreatedCharacter, int, error)
}
