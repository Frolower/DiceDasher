package vtmv5

import (
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
)

type rollState struct {
	system.Lineage
	MainRoll   []int `json:"main_roll"`
	HungerRoll []int `json:"hunger_roll"`
	Target     int   `json:"target"`
}

// loadState adapts legacy initial rolls once; new records carry a complete snapshot.
func loadState(rec repository.RollHistory) (rollState, error) {
	var s rollState
	if err := system.CheckTransition(rec, "vtmv5", "reroll"); err != nil {
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
		s.MainRoll = response.MainRoll
		s.HungerRoll = response.HungerRoll
	}
	if len(s.MainRoll)+len(s.HungerRoll) == 0 || s.Target < 1 {
		return s, errors.New("invalid stored roll state")
	}
	count, err := dice.AddCounts(len(s.MainRoll), len(s.HungerRoll))
	if err != nil {
		return s, err
	}
	if _, err := dice.NewPool(count, dieSize); err != nil {
		return s, err
	}
	for _, pool := range [][]int{s.MainRoll, s.HungerRoll} {
		for _, value := range pool {
			if value < 1 || value > dieSize {
				return s, errors.New("invalid stored die value")
			}
		}
	}
	s.Lineage = s.Lineage.Next(rec.ID)
	return s, nil
}
