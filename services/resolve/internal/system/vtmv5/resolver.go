package vtmv5

import (
	"bytes"
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"fmt"
)

const dieSize = 10

type Resolver struct {
	Dice    dice.Generator
	history system.HistoryReader
}

func New(history system.HistoryReader, generator dice.Generator) Resolver {
	return Resolver{Dice: generator, history: history}
}

func (r Resolver) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, error) {

	switch action {
	case "roll":
		return r.resolveRoll(raw)
	case "reroll":
		return r.resolveReroll(ctx, raw)
	case "check":
		return r.resolveCheck(raw)
	default:
		return nil, system.Invalid(system.BadRequest, errors.New("unknown action"))
	}
}

func (r Resolver) resolveRoll(raw json.RawMessage) (rollResponse, error) {
	var req rollRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return rollResponse{}, system.Invalid(system.BadRequest, err)
	}
	mainPool, hungerPool, err := rollPools(req)
	if err != nil {
		return rollResponse{}, system.Invalid(system.Validation, err)
	}

	expression := fmt.Sprintf("%dd%d", mainPool.Count()+hungerPool.Count(), dieSize)
	mainRoll, err := r.Dice.RollDice(mainPool.Count(), dieSize)
	if err != nil {
		return rollResponse{}, err
	}
	hungerRoll, err := r.Dice.RollDice(hungerPool.Count(), dieSize)
	if err != nil {
		return rollResponse{}, err
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
	}, nil
}

func (r Resolver) resolveReroll(ctx context.Context, raw json.RawMessage) (rerollResponse, error) {
	var req rerollRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return rerollResponse{}, system.Invalid(system.BadRequest, err)
	}
	if err := validateReroll(req); err != nil {
		return rerollResponse{}, system.Invalid(system.Validation, err)
	}

	if r.history == nil {
		return rerollResponse{}, errors.New("history reader is not configured")
	}

	extractedData, err := r.history.GetRollHistoryByID(ctx, req.RecordID)
	if err != nil {
		return rerollResponse{}, err
	}

	rec, err := loadState(extractedData)
	if err != nil {
		return rerollResponse{}, err
	}

	return r.continueRoll(rec, req.RerollIndex)
}

func (r Resolver) continueRoll(rec rollState, indices []int) (rerollResponse, error) {
	if err := validateRerollState(rec, indices); err != nil {
		return rerollResponse{}, system.Invalid(system.Validation, err)
	}

	expression := fmt.Sprintf("%dd%d", len(rec.MainRoll)+len(rec.HungerRoll), dieSize)
	rerollExpression := fmt.Sprintf("%dd%d", len(indices), dieSize)
	mainRoll, err := r.Dice.RerollSpecificValues(rec.MainRoll, indices, dieSize)
	if err != nil {
		return rerollResponse{}, err
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
	}, nil
}

func (r Resolver) resolveCheck(raw json.RawMessage) (checkResponse, error) {
	if len(bytes.TrimSpace(raw)) != 0 {
		return checkResponse{}, system.Invalid(system.BadRequest, errors.New("this action takes an empty body"))
	}

	expression := fmt.Sprintf("1d%d", dieSize)
	result, err := r.Dice.RollDie(dieSize)
	if err != nil {
		return checkResponse{}, err
	}
	success := result >= 6

	return checkResponse{
		Expression: expression,
		Result:     result,
		Success:    success,
	}, nil
}
