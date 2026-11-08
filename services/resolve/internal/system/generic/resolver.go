package generic

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/util"
	"diceDasher/services/resolve/internal/system"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
)

type Resolver struct{ Dice dice.Generator }

func (r Resolver) Resolve(ctx context.Context, action string, raw jsontext.Value) (any, error) {
	var req request

	if action != "roll" {
		return nil, system.Invalid(system.BadRequest, errors.New("invalid action"))
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, system.Invalid(system.BadRequest, err)
	}
	if err := validate(req); err != nil {
		return nil, system.Invalid(system.Validation, err)
	}

	expression := fmt.Sprintf("%dd%d", req.Number, req.Size)
	rolls, err := r.Dice.RollDice(req.Number, req.Size)
	if err != nil {
		return response{}, err
	}
	sum := util.Sum(rolls)

	return response{
		Expression: expression,
		Rolls:      rolls,
		Sum:        sum,
	}, nil
}
