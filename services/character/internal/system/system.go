package system

import (
	"context"
	"encoding/json"
)

type Character interface {
	CreateCharacter(ctx context.Context, raw json.RawMessage) (any, int, error)
}
