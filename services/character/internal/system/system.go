package system

import (
	"context"
	"encoding/json/jsontext"

	"github.com/google/uuid"
)

type CreatedCharacter struct {
	UserID        uuid.UUID
	CharacterType string
	Name          string
	Data          jsontext.Value
}

type Character interface {
	CreateCharacter(ctx context.Context, raw jsontext.Value) (CreatedCharacter, error)
}
