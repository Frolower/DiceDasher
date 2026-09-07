package tes

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/util"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"fmt"
)

const dieSize = 6

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
	case "push":
		return r.resolvePush(ctx, raw)
	default:
		return nil, system.Invalid(system.BadRequest, errors.New("invalid action"))
	}
}

func (r Resolver) resolveRoll(raw json.RawMessage) (rollResponse, error) {
	var req rollRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return rollResponse{}, system.Invalid(system.BadRequest, err)
	}
	attributePool, gearPool, err := rollPools(req)
	if err != nil {
		return rollResponse{}, system.Invalid(system.Validation, err)
	}

	expression := fmt.Sprintf("%dd%d", attributePool.Count()+gearPool.Count(), dieSize)
	attributeRolls, err := r.Dice.RollDice(attributePool.Count(), dieSize)
	if err != nil {
		return rollResponse{}, err
	}
	gearRolls, err := r.Dice.RollDice(gearPool.Count(), dieSize)
	if err != nil {
		return rollResponse{}, err
	}
	successes := util.CountInt(attributeRolls, 6) + util.CountInt(gearRolls, 6)
	success := successes >= req.Target

	return rollResponse{
		state:          rollState{Target: req.Target, AttributeRolls: attributeRolls, GearRolls: gearRolls},
		Expression:     expression,
		AttributeRolls: attributeRolls,
		GearRolls:      gearRolls,
		Successes:      successes,
		Success:        success,
	}, nil
}

func (r Resolver) resolvePush(ctx context.Context, raw json.RawMessage) (pushResponse, error) {
	var req pushRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return pushResponse{}, system.Invalid(system.BadRequest, err)
	}
	if err := validatePush(req); err != nil {
		return pushResponse{}, system.Invalid(system.Validation, err)
	}

	if r.history == nil {
		return pushResponse{}, errors.New("history reader is not configured")
	}

	extractedData, err := r.history.GetRollHistoryByID(ctx, req.RecordID)
	if err != nil {
		return pushResponse{}, err
	}

	rec, err := loadState(extractedData)
	if err != nil {
		return pushResponse{}, err
	}

	return r.continueRoll(rec)
}

func (r Resolver) continueRoll(rec rollState) (pushResponse, error) {
	expression := fmt.Sprintf("%dd%d", len(rec.AttributeRolls)+len(rec.GearRolls), dieSize)
	rerollDiceNumber := util.CountBetween(rec.AttributeRolls, 2, 5) + util.CountBetween(rec.GearRolls, 2, 5)
	pushExpression := fmt.Sprintf("%dd%d", rerollDiceNumber, dieSize)
	attributeRolls, err := r.Dice.RerollKeepingValues(rec.AttributeRolls, []int{1, 6}, dieSize)
	if err != nil {
		return pushResponse{}, err
	}
	gearRolls, err := r.Dice.RerollKeepingValues(rec.GearRolls, []int{1, 6}, dieSize)
	if err != nil {
		return pushResponse{}, err
	}
	successes := util.CountInt(attributeRolls, 6) + util.CountInt(gearRolls, 6)
	success := successes >= rec.Target
	hopeLosses := util.CountInt(attributeRolls, 1)
	gearDamage := util.CountInt(gearRolls, 1)

	return pushResponse{
		state:          rollState{Target: rec.Target, AttributeRolls: attributeRolls, GearRolls: gearRolls, Lineage: rec.Lineage},
		Expression:     expression,
		PushExpression: pushExpression,
		AttributeRolls: attributeRolls,
		GearRolls:      gearRolls,
		Successes:      successes,
		Success:        success,
		HopeLosses:     hopeLosses,
		GearDamage:     gearDamage,
	}, nil
}
