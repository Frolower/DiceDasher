package tes

import (
	"context"
	"diceDasher/pkg/logger"
	"encoding/json"
	"net/http"
)

type Character struct{}

func (Character) CreateCharacter(ctx context.Context, raw json.RawMessage) (any, int, error) {
	logger.Logf(ctx, "RUN: CreateCharacter system=tes |")
	var req createRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err := validateCreate(req); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

}
