package tes

import (
	"context"
	"diceDasher/pkg/logger"
	"diceDasher/services/character/internal/system"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Character struct{}

func (Character) CreateCharacter(ctx context.Context, raw json.RawMessage) (system.CreatedCharacter, int, error) {
	logger.Logf(ctx, "RUN: CreateCharacter system=tes |")

	var req createRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return system.CreatedCharacter{}, http.StatusBadRequest, err
	}
	if req.UserID == uuid.Nil {
		return system.CreatedCharacter{}, http.StatusBadRequest, errors.New("user_id is required")
	}
	if err := validateCreate(req); err != nil {
		return system.CreatedCharacter{}, http.StatusUnprocessableEntity, err
	}

	data, err := json.Marshal(storedCharacter{
		Type:      req.Type,
		Character: req.CharacterList,
	})
	if err != nil {
		return system.CreatedCharacter{}, http.StatusInternalServerError, err
	}

	return system.CreatedCharacter{
		UserID:        req.UserID,
		CharacterType: req.Type,
		Name:          req.CharacterList.Name,
		Data:          data,
	}, http.StatusCreated, nil
}
