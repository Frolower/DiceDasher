package tes

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/logger"
	"diceDasher/pkg/util"
	"diceDasher/services/resolve/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const dieSize = 6

type Resolver struct{}

func (Resolver) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error) {
	logger.Logf(ctx, "RUN: resolver=tes action=%s |", action)

	switch action {
	case "roll":
		return resolveRoll(raw)
	case "push":
		return resolvePush(ctx, raw)
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

	expression := fmt.Sprintf("%dd%d", req.Attr+req.Assist+req.Gear+req.Modificator, dieSize)
	attributeRolls, err := dice.RollDice(req.Attr+req.Assist+req.Modificator, dieSize)
	if err != nil {
		return rollResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	gearRolls, err := dice.RollDice(req.Gear, dieSize)
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

func resolvePush(ctx context.Context, raw json.RawMessage) (pushResponse, int, error) {
	var req pushRequest
	var rec pushRecord

	if err := json.Unmarshal(raw, &req); err != nil {
		return pushResponse{}, http.StatusBadRequest, err
	}
	if err := validatePush(req); err != nil {
		return pushResponse{}, http.StatusUnprocessableEntity, err
	}

	repo, err := repository.FromContext(ctx)
	if err != nil {
		logger.Logf(ctx, "ERROR: repository not in context: %s", err)
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}

	extractedData, err := repo.GetRollHistoryByID(ctx, req.RecordID)
	if err != nil {
		return pushResponse{}, http.StatusInternalServerError, err
	}

	// Extracting results from thr db
	if err = json.Unmarshal(extractedData.ResponsePayload, &rec); err != nil {
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}

	// Extracting target from the request
	if err = json.Unmarshal(extractedData.RequestPayload, &rec); err != nil {
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}

	if err := validatePushRecord(rec); err != nil {
		return pushResponse{}, http.StatusInternalServerError, err
	}

	expression := fmt.Sprintf("%dd%d", len(rec.AttributeRolls)+len(rec.GearRolls), dieSize)
	rerollDiceNumber := util.CountBetween(rec.AttributeRolls, 2, 5) + util.CountBetween(rec.GearRolls, 2, 5)
	pushExpression := fmt.Sprintf("%dd%d", rerollDiceNumber, dieSize)
	attributeRolls, err := dice.RerollKeepingValues(rec.AttributeRolls, []int{1, 6}, dieSize)
	if err != nil {
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	gearRolls, err := dice.RerollKeepingValues(rec.GearRolls, []int{1, 6}, dieSize)
	if err != nil {
		return pushResponse{}, http.StatusInternalServerError, errors.New("internal error")
	}
	successes := util.CountInt(attributeRolls, 6) + util.CountInt(gearRolls, 6)
	success := successes >= rec.Target
	hopeLosses := util.CountInt(attributeRolls, 1)
	gearDamage := util.CountInt(gearRolls, 1)

	return pushResponse{
		Expression:     expression,
		PushExpression: pushExpression,
		AttributeRolls: attributeRolls,
		GearRolls:      gearRolls,
		Successes:      successes,
		Success:        success,
		HopeLosses:     hopeLosses,
		GearDamage:     gearDamage,
	}, http.StatusOK, nil
}
