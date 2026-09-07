package tes

import (
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
)

type rollState struct {
	system.Lineage
	AttributeRolls []int `json:"attribute_rolls"`
	GearRolls      []int `json:"gear_rolls"`
	Target         int   `json:"target"`
}

// loadState adapts legacy initial rolls once; new records carry a complete snapshot.
func loadState(rec repository.RollHistory) (rollState, error) {
	var s rollState
	if err := system.CheckTransition(rec, "tes", "push"); err != nil {
		return s, err
	}
	if system.HasState(rec.StatePayload) {
		if err := json.Unmarshal(rec.StatePayload, &s); err != nil {
			return s, err
		}
	} else {
		if rec.ActionType != "roll" {
			return s, system.ErrLegacyContinuation
		}
		var request rollRequest
		var response rollResponse
		if err := json.Unmarshal(rec.RequestPayload, &request); err != nil {
			return s, err
		}
		if err := json.Unmarshal(rec.ResponsePayload, &response); err != nil {
			return s, err
		}
		s.Target = request.Target
		s.AttributeRolls = response.AttributeRolls
		s.GearRolls = response.GearRolls
	}
	if len(s.AttributeRolls)+len(s.GearRolls) == 0 || s.Target < 0 {
		return s, errors.New("invalid stored roll state")
	}
	for _, pool := range [][]int{s.AttributeRolls, s.GearRolls} {
		for _, value := range pool {
			if value < 1 || value > dieSize {
				return s, errors.New("invalid stored die value")
			}
		}
	}
	s.Lineage = s.Lineage.Next(rec.ID)
	return s, nil
}
