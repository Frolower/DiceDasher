package system

import (
	"diceDasher/services/resolve/internal/repository"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

// StatefulResult exposes persistence state without adding fields to the HTTP payload.
type StatefulResult interface{ HistoryState() any }

// Lineage records the original roll and the immediate predecessor of a continuation.
// A new roll has no predecessor; its original ID is assigned when it is loaded.
type Lineage struct {
	OriginalID uuid.UUID `json:"original_id"`
	ParentID   uuid.UUID `json:"parent_id"`
}

var ErrInvalidTransition = errors.New("record does not support this action")
var ErrLegacyContinuation = errors.New("legacy continuation has no complete state; use the original roll")

func CheckTransition(rec repository.RollHistory, name, continuation string) error {
	if rec.SystemName != name || (rec.ActionType != "roll" && rec.ActionType != continuation) {
		return ErrInvalidTransition
	}
	return nil
}

func (s Lineage) Next(previous uuid.UUID) Lineage {
	if s.OriginalID == uuid.Nil {
		s.OriginalID = previous
	}
	s.ParentID = previous
	return s
}

func HasState(raw json.RawMessage) bool { return len(raw) != 0 && string(raw) != "null" }

func HistoryErrorStatus(err error) int {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidTransition), errors.Is(err, ErrLegacyContinuation):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
