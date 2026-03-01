package tes

import (
	"context"
	"diceDasher/pkg/dice"
	"encoding/json"
	"fmt"
	"net/http"
)

type Resolver struct{}

func (Resolver) Resolve(ctx context.Context, raw json.RawMessage) (any, int, error) {
	var req request
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err := validate(req); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd6", req.Attr+req.Assist+req.Gear+req.Modificator)
	rolls := dice.RollDice(req.Attr+req.Assist+req.Gear+req.Modificator, 6)
	successes := dice.CountInt(rolls, 6)
	success := successes >= req.Target

	return response{
		Expression: expression,
		Rolls:      rolls,
		Successes:  successes,
		Success:    success,
	}, http.StatusOK, nil
}
