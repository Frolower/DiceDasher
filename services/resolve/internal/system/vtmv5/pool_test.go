package vtmv5

import (
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"testing"
)

func TestRollPoolValidation(t *testing.T) {
	max := int(^uint(0) >> 1)
	for _, req := range []rollRequest{
		{Attribute: max, Skill: 1, Target: 1},
		{Attribute: dice.MaxDice, Skill: 1, Target: 1},
		{Attribute: 1, Hunger: 2, Target: 1},
		{Attribute: 1, Hunger: -1, Target: 1},
	} {
		raw, _ := json.Marshal(req)
		_, err := (Resolver{}).resolveRoll(raw)
		if !system.IsValidation(err) {
			t.Fatalf("accepted %+v: %v", req, err)
		}
	}
	for _, req := range []rollRequest{
		{Attribute: dice.MaxDice, Target: 1},
		{Attribute: 2, Hunger: 2, Target: 1},
	} {
		raw, _ := json.Marshal(req)
		result, err := (Resolver{}).resolveRoll(raw)
		if err != nil {
			t.Fatalf("rejected %+v: %v", req, err)
		}
		if len(result.MainRoll)+len(result.HungerRoll) != req.Attribute+req.Skill {
			t.Fatal("wrong pool size")
		}
	}
}

func TestStoredPoolLimit(t *testing.T) {
	values := make([]int, dice.MaxDice)
	for i := range values {
		values[i] = 1
	}
	state := rollState{MainRoll: values, HungerRoll: []int{1}, Target: 1}
	raw, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := loadState(repository.RollHistory{SystemName: "vtmv5", ActionType: "roll", StatePayload: raw}); err == nil {
		t.Fatal("oversized combined state accepted")
	}
}
