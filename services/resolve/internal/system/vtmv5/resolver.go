package vtmv5

import (
	"bytes"
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/logger"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const dieSize = 10

type Resolver struct{ Dice dice.Generator }

func (r Resolver) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error) {
	logger.Logf(ctx, "RUN: resolver=vtmv5 action=%s |", action)

	switch action {
	case "roll":
		return r.resolveRoll(raw)
	case "reroll":
		return r.resolveReroll(ctx, raw)
	case "check":
		return r.resolveCheck(raw)
	default:
		return nil, http.StatusBadRequest, errors.New("unknown action")
	}
}

func (r Resolver) resolveRoll(raw json.RawMessage) (rollResponse, int, error) {
	var req rollRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return rollResponse{}, http.StatusBadRequest, err
	}
	mainPool, hungerPool, err := rollPools(req)
	if err != nil {
		return rollResponse{}, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd%d", mainPool.Count()+hungerPool.Count(), dieSize)
	mainRoll, err := r.Dice.RollDice(mainPool.Count(), dieSize)
	if err != nil {
		return rollResponse{}, http.StatusBadRequest, errors.New("internal error")
	}
	hungerRoll, err := r.Dice.RollDice(hungerPool.Count(), dieSize)
	if err != nil {
		return rollResponse{}, http.StatusBadRequest, errors.New("internal error")
	}
	outcome := Evaluate(mainRoll, hungerRoll, req.Target)

	return rollResponse{
		state:      rollState{Target: req.Target, MainRoll: mainRoll, HungerRoll: hungerRoll},
		Expression: expression,
		MainRoll:   mainRoll,
		HungerRoll: hungerRoll,
		Successes:  outcome.Successes,
		Success:    outcome.Success,
		IsCritical: outcome.IsCritical,
		CritType:   outcome.CritType,
	}, http.StatusOK, nil
}

func (r Resolver) resolveReroll(ctx context.Context, raw json.RawMessage) (rerollResponse, int, error) {
	var req rerollRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return rerollResponse{}, http.StatusBadRequest, err
	}
	if err := validateReroll(req); err != nil {
		return rerollResponse{}, http.StatusUnprocessableEntity, err
	}

	repo, err := repository.FromContext(ctx)
	if err != nil {
		logger.Logf(ctx, "ERROR: repository not in context: %s", err)
		return rerollResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}

	extractedData, err := repo.GetRollHistoryByID(ctx, req.RecordID)
	if err != nil {
		return rerollResponse{}, system.HistoryErrorStatus(err), err
	}

	rec, err := loadState(extractedData)
	if err != nil {
		return rerollResponse{}, system.HistoryErrorStatus(err), err
	}

	return r.continueRoll(rec, req.RerollIndex)
}

func (r Resolver) continueRoll(rec rollState, indices []int) (rerollResponse, int, error) {
	if err := validateRerollState(rec, indices); err != nil {
		return rerollResponse{}, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd%d", len(rec.MainRoll)+len(rec.HungerRoll), dieSize)
	rerollExpression := fmt.Sprintf("%dd%d", len(indices), dieSize)
	mainRoll, err := r.Dice.RerollSpecificValues(rec.MainRoll, indices, dieSize)
	if err != nil {
		return rerollResponse{}, http.StatusBadRequest, errors.New("interal error")
	}
	hungerRoll := rec.HungerRoll
	outcome := Evaluate(mainRoll, hungerRoll, rec.Target)

	return rerollResponse{
		state:            rollState{Target: rec.Target, MainRoll: mainRoll, HungerRoll: hungerRoll, Lineage: rec.Lineage},
		Expression:       expression,
		RerollExpression: rerollExpression,
		MainRoll:         mainRoll,
		HungerRoll:       hungerRoll,
		Successes:        outcome.Successes,
		Success:          outcome.Success,
		IsCritical:       outcome.IsCritical,
		CritType:         outcome.CritType,
	}, http.StatusOK, nil
}

func (r Resolver) resolveCheck(raw json.RawMessage) (checkResponse, int, error) {
	if len(bytes.TrimSpace(raw)) != 0 {
		return checkResponse{}, http.StatusBadRequest, errors.New("this action takes an empty body")
	}

	expression := fmt.Sprintf("1d%d", dieSize)
	result, err := r.Dice.RollDie(dieSize)
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
