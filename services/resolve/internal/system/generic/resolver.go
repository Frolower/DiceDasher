package generic

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

	expression := fmt.Sprintf("%dd%d", req.Number, req.Size)
	rolls := dice.RollDice(req.Number, req.Size)
	sum := dice.Sum(rolls)

	return response{
		Expression: expression,
		Rolls:      rolls,
		Sum:        sum,
	}, http.StatusOK, nil
}
