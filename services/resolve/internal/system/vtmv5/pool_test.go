package vtmv5

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
		{Attribute: max, Skill: 1, Target: 1},
		{Attribute: dice.MaxDice, Skill: 1, Target: 1},
		{Attribute: 1, Hunger: 2, Target: 1},
		{Attribute: 1, Hunger: -1, Target: 1},
	} {
		raw, _ := json.Marshal(req)
		_, status, err := (Resolver{}).resolveRoll(raw)
		if err == nil || status != http.StatusUnprocessableEntity {
			t.Fatalf("accepted %+v: %d %v", req, status, err)
		}
	}
	for _, req := range []rollRequest{
		{Attribute: dice.MaxDice, Target: 1},
		{Attribute: 2, Hunger: 2, Target: 1},
	} {
		raw, _ := json.Marshal(req)
		result, status, err := (Resolver{}).resolveRoll(raw)
		if err != nil || status != http.StatusOK {
			t.Fatalf("rejected %+v: %d %v", req, status, err)
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
