package tes

import (
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/repository"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRollPoolValidation(t *testing.T) {
	max := int(^uint(0) >> 1)
	for _, req := range []rollRequest{
		{Attr: 1, Gear: 2, Modificator: -2, Target: 1},
		{Attr: max, Assist: 1},
		{Attr: 1, Modificator: max},
		{Attr: 1, Modificator: -max},
		{Attr: dice.MaxDice, Gear: 1},
	} {
		raw, _ := json.Marshal(req)
		_, status, err := resolveRoll(raw)
		if err == nil || status != http.StatusUnprocessableEntity {
			t.Fatalf("accepted %+v: %d %v", req, status, err)
		}
	}
	for _, req := range []rollRequest{
		{Attr: 1, Gear: 2, Modificator: -1, Target: 2}, // empty attribute component is valid
		{Attr: dice.MaxDice, Target: 0},
		{Attr: 1, Modificator: dice.MaxDice - 1, Target: 1},
	} {
		raw, _ := json.Marshal(req)
		result, status, err := resolveRoll(raw)
		if err != nil || status != http.StatusOK {
			t.Fatalf("rejected %+v: %d %v", req, status, err)
		}
		if len(result.AttributeRolls)+len(result.GearRolls) != req.Attr+req.Assist+req.Gear+req.Modificator {
			t.Fatal("wrong pool size")
		}
	}
}

func TestStoredPoolLimit(t *testing.T) {
	values := make([]int, dice.MaxDice)
	for i := range values {
		values[i] = 1
	}
	state := rollState{AttributeRolls: values, GearRolls: []int{1}, Target: 1}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := loadState(repository.RollHistory{SystemName: "tes", ActionType: "roll", StatePayload: raw}); err == nil {
		t.Fatal("oversized combined state accepted")
	}
}
