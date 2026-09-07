package tes

import (
	"context"
	"diceDasher/pkg/logger"
	"diceDasher/services/character/internal/system"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type Character struct{}

func (Character) CreateCharacter(ctx context.Context, raw json.RawMessage) (system.CreatedCharacter, error) {
	logger.Logf(ctx, "RUN: CreateCharacter system=tes |")

	var req createRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return system.CreatedCharacter{}, &system.InputError{Err: err}
	}
	if req.UserID == uuid.Nil {
		return system.CreatedCharacter{}, &system.InputError{Err: errors.New("user_id is required")}
	}
	if err := validateCreate(req); err != nil {
		return system.CreatedCharacter{}, err
	}

	data, err := json.Marshal(storedCharacter{
		Type:      req.Type,
		Character: req.CharacterList,
	})
	if err != nil {
		return system.CreatedCharacter{}, err
	}

	return system.CreatedCharacter{
		UserID:        req.UserID,
		CharacterType: req.Type,
		Name:          req.CharacterList.Name,
		Data:          data,
	}, nil
}
