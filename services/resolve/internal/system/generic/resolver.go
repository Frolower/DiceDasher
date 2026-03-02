package generic

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Resolver struct{}

func (Resolver) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error) {
	httputil.Logf(ctx, "RUN: resolver=generic action=%s", action)
	var req request

	if action != "roll" {
		return nil, http.StatusBadRequest, errors.New("invalid action")
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err := validate(req); err != nil {
		return nil, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd%d", req.Number, req.Size)
	rolls, err := dice.RollDice(req.Number, req.Size)
	if err != nil {
		return response{}, http.StatusInternalServerError, errors.New("internal error")
	}
	sum := util.Sum(rolls)

	return response{
		Expression: expression,
		Rolls:      rolls,
		Sum:        sum,
	}, http.StatusOK, nil
}
