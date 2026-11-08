package system

import (
	"context"
	"diceDasher/services/resolve/internal/repository"
	"encoding/json/jsontext"
	"github.com/google/uuid"
)

type Resolver interface {
	Resolve(ctx context.Context, action string, raw jsontext.Value) (any, error)
}

// HistoryReader is the only persistence capability needed by continuations.
type HistoryReader interface {
	GetRollHistoryByID(context.Context, uuid.UUID) (repository.RollHistory, error)
}
