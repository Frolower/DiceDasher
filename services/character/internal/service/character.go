package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"

	"diceDasher/services/character/internal/system"
	"github.com/google/uuid"
)

var ErrUnknownSystem = errors.New("unknown system")

// Repository is the persistence capability required by the creation use case.
type Repository interface {
	InsertPlayerCreatedCharacter(context.Context, Record) (uuid.UUID, error)
}

type Record struct {
	UserID        uuid.UUID
	SystemName    string
	CharacterType string
	Name          string
	Data          json.RawMessage
}

type SaveError struct{ Err error }

func (e *SaveError) Error() string { return "failed to save character: " + e.Err.Error() }
func (e *SaveError) Unwrap() error { return e.Err }

type CharacterService struct {
	repository Repository
	systems    map[string]system.Character
}

// New snapshots the configured strategies so later map changes cannot alter the service.
func New(repository Repository, systems map[string]system.Character) *CharacterService {
	return &CharacterService{repository: repository, systems: maps.Clone(systems)}
}

func (s *CharacterService) Create(ctx context.Context, systemName string, raw json.RawMessage) (uuid.UUID, error) {
	strategy, ok := s.systems[systemName]
	if !ok {
		return uuid.Nil, ErrUnknownSystem
	}
	created, err := strategy.CreateCharacter(ctx, raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create character: %w", err)
	}
	id, err := s.repository.InsertPlayerCreatedCharacter(ctx, Record{
		UserID: created.UserID, SystemName: systemName, CharacterType: created.CharacterType, Name: created.Name, Data: created.Data,
	})
	if err != nil {
		return uuid.Nil, &SaveError{Err: err}
	}
	return id, nil
}
