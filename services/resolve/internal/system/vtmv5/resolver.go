package vtmv5

import (
	"bytes"
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const dieSize = 10

type Resolver struct{}

func (Resolver) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error) {
	switch action {
	case "roll":
		return resolveRoll(raw)
	case "reroll":
		return resolveReroll(raw)
	case "check":
		return resolveCheck(raw)
	default:
		return nil, http.StatusBadRequest, errors.New("unknown action")
	}
}

func resolveRoll(raw json.RawMessage) (rollResponse, int, error) {
	req := rollRequest{}

	if err := json.Unmarshal(raw, &req); err != nil {
		return rollResponse{}, http.StatusBadRequest, err
	}
	if err := validateRoll(req); err != nil {
		return rollResponse{}, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd%d", req.Attribute+req.Skill, dieSize)
	mainRoll, err := dice.RollDice(util.Max(req.Attribute+req.Skill-req.Hunger, 0), dieSize)
	if err != nil {
		return rollResponse{}, http.StatusBadRequest, errors.New("internal error")
	}
	hungerRoll, err := dice.RollDice(req.Hunger, dieSize)
	if err != nil {
		return rollResponse{}, http.StatusBadRequest, errors.New("internal error")
	}
	regular10s := util.CountInt(mainRoll, 10)
	hunger10s := util.CountInt(hungerRoll, 10)
	successes := util.CountBetween(mainRoll, 6, 10) + util.CountBetween(hungerRoll, 6, 10)
	pairsOf10s := (regular10s + hunger10s) / 2
	successes += pairsOf10s * 2 // Adds 2 extra successes for each pair of 10's
	success := successes >= req.Target
	hungerFails := util.CountInt(hungerRoll, 1)
	isCritical := false
	critType := "none"

	if !success && hungerFails > 0 {
		isCritical = true
		critType = "bestial failure"
	}

	if pairsOf10s > 0 {
		isCritical = true
		if hunger10s > 0 {
			critType = "messy critical"
		} else {
			critType = "critical"
		}
	}

	return rollResponse{
		Expression: expression,
		MainRoll:   mainRoll,
		HungerRoll: hungerRoll,
		Successes:  successes,
		Success:    success,
		IsCritical: isCritical,
		CritType:   critType,
	}, http.StatusOK, nil
}

func resolveReroll(raw json.RawMessage) (rerollResponse, int, error) {
	req := rerollRequest{}

	if err := json.Unmarshal(raw, &req); err != nil {
		return rerollResponse{}, http.StatusBadRequest, err
	}
	if err := validateReroll(req); err != nil {
		return rerollResponse{}, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd%d", len(req.MainRoll)+len(req.HungerRoll), dieSize)
	rerollExpression := fmt.Sprintf("%dd%d", len(req.RerollIndex), dieSize)
	mainRoll, err := dice.RerollSpecificValues(req.MainRoll, req.RerollIndex, dieSize)
	if err != nil {
		return rerollResponse{}, http.StatusBadRequest, errors.New("interal error")
	}
	hungerRoll := req.HungerRoll
	regular10s := util.CountInt(mainRoll, 10)
	hunger10s := util.CountInt(hungerRoll, 10)
	successes := util.CountBetween(mainRoll, 6, 10) + util.CountBetween(hungerRoll, 6, 10)
	pairsOf10s := (regular10s + hunger10s) / 2
	successes += pairsOf10s * 2 // Adds 2 extra successes for each pair of 10's
	success := successes >= req.Target
	hungerFails := util.CountInt(hungerRoll, 1)
	isCritical := false
	critType := "none"

	if !success && hungerFails > 0 {
		isCritical = true
		critType = "bestial failure"
	}

	if pairsOf10s > 0 {
		isCritical = true
		if hunger10s > 0 {
			critType = "messy critical"
		} else {
			critType = "critical"
		}
	}

	return rerollResponse{
		Expression:       expression,
		RerollExpression: rerollExpression,
		MainRoll:         mainRoll,
		HungerRoll:       hungerRoll,
		Successes:        successes,
		Success:          success,
		IsCritical:       isCritical,
		CritType:         critType,
	}, http.StatusOK, nil
}

func resolveCheck(raw json.RawMessage) (checkResponse, int, error) {
	if len(bytes.TrimSpace(raw)) != 0 {
		return checkResponse{}, http.StatusBadRequest, errors.New("this action takes an empty body")
	}

	expression := fmt.Sprintf("1d%d", dieSize)
	result, err := dice.RollDie(dieSize)
	if err != nil {
		return checkResponse{}, http.StatusUnprocessableEntity, errors.New("internal error")
	}
	success := result >= 6

	return checkResponse{
		Expression: expression,
		Result:     result,
		Success:    success,
	}, http.StatusOK, nil
}
