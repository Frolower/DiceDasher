package tes

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Resolver struct{}

func (Resolver) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error) {
	switch action {
	case "roll":
		return resolveRoll(raw)
	case "push":
		return resolvePush(raw)
	default:
		return nil, http.StatusBadRequest, errors.New("invalid action")
	}
}

func resolveRoll(raw json.RawMessage) (rollResponse, int, error) {
	var req rollRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return rollResponse{}, http.StatusBadRequest, err
	}
	if err := validateRoll(req); err != nil {
		return rollResponse{}, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd6", req.Attr+req.Assist+req.Gear+req.Modificator)
	attributeRolls, err := dice.RollDice(req.Attr+req.Assist+req.Modificator, 6)
	if err != nil {
		return rollResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	gearRolls, err := dice.RollDice(req.Gear, 6)
	if err != nil {
		return rollResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	successes := util.CountInt(attributeRolls, 6) + util.CountInt(gearRolls, 6)
	success := successes >= req.Target

	return rollResponse{
		Expression:     expression,
		AttributeRolls: attributeRolls,
		GearRolls:      gearRolls,
		Successes:      successes,
		Success:        success,
	}, http.StatusOK, nil
}

func resolvePush(raw json.RawMessage) (pushResponse, int, error) {
	var req pushRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		return pushResponse{}, http.StatusBadRequest, err
	}
	if err := validatePush(req); err != nil {
		return pushResponse{}, http.StatusUnprocessableEntity, err
	}

	expression := fmt.Sprintf("%dd6", len(req.AttributeRolls)+len(req.GearRolls))
	attributeRolls, err := dice.RerollKeepingValues(req.AttributeRolls, []int{1, 6}, 6)
	if err != nil {
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	gearRolls, err := dice.RerollKeepingValues(req.GearRolls, []int{1, 6}, 6)
	if err != nil {
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	successes := util.CountInt(attributeRolls, 6) + util.CountInt(gearRolls, 6)
	success := successes >= req.Target
	hopeLosses := util.CountInt(attributeRolls, 1)
	gearDamage := util.CountInt(gearRolls, 1)

	return pushResponse{
		Expression:     expression,
		AttributeRolls: attributeRolls,
		GearRolls:      gearRolls,
		Successes:      successes,
		Success:        success,
		HopeLosses:     hopeLosses,
		GearDamage:     gearDamage,
	}, http.StatusOK, nil
}
