package service

import (
	"context"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type HistoryWriter interface {
	InsertRollHistory(context.Context, repository.RollHistory) (uuid.UUID, error)
}
type ResolveService struct {
	resolvers map[string]system.Resolver
	history   HistoryWriter
}

func New(resolvers map[string]system.Resolver, history HistoryWriter) *ResolveService {
	copy := make(map[string]system.Resolver, len(resolvers))
	for name, resolver := range resolvers {
		copy[name] = resolver
	}
	return &ResolveService{resolvers: copy, history: history}
}

type Command struct {
	System    string
	Action    string
	Payload   json.RawMessage
	RequestID uuid.UUID
}
type Result struct {
	RecordID uuid.UUID       `json:"record_id"`
	Payload  json.RawMessage `json:"payload"`
}

// Resolve returns success only after the result and its continuation state have
// been stored together. It never reruns a random operation after a storage error.
func (s *ResolveService) Resolve(ctx context.Context, cmd Command) (Result, error) {
	if cmd.System == "" {
		return Result{}, system.Invalid(system.BadRequest, errors.New("missing query param: system"))
	}
	if cmd.Action == "" {
		cmd.Action = "roll"
	}
	resolver, ok := s.resolvers[cmd.System]
	if !ok {
		return Result{}, system.ErrUnknownSystem
	}
	if s.history == nil {
		return Result{}, errors.New("history writer is not configured")
	}
	if cmd.RequestID == uuid.Nil {
		cmd.RequestID = uuid.New()
	}
	response, err := resolver.Resolve(ctx, cmd.Action, cmd.Payload)
	if err != nil {
		return Result{}, err
	}
	payload, err := json.Marshal(response)
	if err != nil {
		return Result{}, fmt.Errorf("encode result: %w", err)
	}
	var state json.RawMessage
	if stateful, ok := response.(system.StatefulResult); ok {
		state, err = json.Marshal(stateful.HistoryState())
		if err != nil {
			return Result{}, fmt.Errorf("encode state: %w", err)
		}
	}
	id, err := s.history.InsertRollHistory(ctx, repository.RollHistory{
		RequestID: cmd.RequestID, SystemName: cmd.System, ActionType: cmd.Action,
		RequestPayload: cmd.Payload, ResponsePayload: payload, StatePayload: state,
	})
	if err != nil {
		return Result{}, fmt.Errorf("save roll history: %w", err)
	}
	return Result{RecordID: id, Payload: payload}, nil
}
