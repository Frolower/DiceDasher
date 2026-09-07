package vtmv5

import (
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestContinuationRoundTrip(t *testing.T) {
	root := uuid.New()
	rec := repository.RollHistory{ID: root, SystemName: "vtmv5", ActionType: "roll",
		RequestPayload:  json.RawMessage(`{"target":2}`),
		ResponsePayload: json.RawMessage(`{"main_roll":[1,6],"hunger_roll":[1,6]}`),
	}
	for i := 0; i < 3; i++ {
		state, err := loadState(rec)
		if err != nil {
			t.Fatal(err)
		}
		if state.Target != 2 || state.OriginalID != root || state.ParentID != rec.ID {
			t.Fatalf("lost state: %+v", state)
		}
		response, err := (Resolver{}).continueRoll(state, []int{0})
		if err != nil {
			t.Fatalf("continuation: %v", err)
		}

		if response.RerollExpression != "1d10" {
			t.Fatalf("current indices ignored: %s", response.RerollExpression)
		}
		if response.MainRoll[1] != 6 || !reflect.DeepEqual(response.HungerRoll, []int{1, 6}) {
			t.Fatal("unselected dice changed")
		}

		if !reflect.DeepEqual(state.HungerRoll, []int{1, 6}) {
			t.Fatal("source state mutated")
		}
		encoded, err := json.Marshal(response.HistoryState())
		if err != nil {
			t.Fatal(err)
		}
		payload, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}
		var public map[string]json.RawMessage
		if err := json.Unmarshal(payload, &public); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"state", "target", "original_id", "parent_id"} {
			if _, ok := public[key]; ok {
				t.Fatalf("private state exposed: %s", key)
			}
		}
		rec = repository.RollHistory{ID: uuid.New(), SystemName: "vtmv5", ActionType: "reroll", StatePayload: encoded,
			RequestPayload: json.RawMessage(`{"record_id":"` + rec.ID.String() + `"}`), ResponsePayload: payload}
	}
}

func TestRejectUnsupportedHistory(t *testing.T) {
	for _, tc := range []struct {
		system, action string
		expected       error
	}{
		{"generic", "roll", system.ErrInvalidTransition},
		{"vtmv5", "check", system.ErrInvalidTransition},
		{"vtmv5", "reroll", system.ErrLegacyContinuation},
	} {
		_, err := loadState(repository.RollHistory{SystemName: tc.system, ActionType: tc.action})
		if !errors.Is(err, tc.expected) {
			t.Fatalf("%s/%s: %v", tc.system, tc.action, err)
		}
	}
}

func TestInitialRollPersistsState(t *testing.T) {
	raw := json.RawMessage(`{"attribute":2,"skill":1,"hunger":1,"target":1}`)
	response, err := (Resolver{}).resolveRoll(raw)
	if err != nil {
		t.Fatalf("roll: %v", err)
	}
	encoded, err := json.Marshal(response.HistoryState())
	if err != nil {
		t.Fatal(err)
	}
	state, err := loadState(repository.RollHistory{ID: uuid.New(), SystemName: "vtmv5", ActionType: "roll", StatePayload: encoded})
	if err != nil {
		t.Fatal(err)
	}
	if state.Target != 1 || !reflect.DeepEqual(state.MainRoll, response.MainRoll) || !reflect.DeepEqual(state.HungerRoll, response.HungerRoll) {
		t.Fatal("snapshot differs from result")
	}
}

func TestInvalidRerollIndices(t *testing.T) {
	state := rollState{MainRoll: []int{1, 6}, Target: 1}
	for _, indices := range [][]int{{-1}, {2}, {0, 0}} {
		_, err := (Resolver{}).continueRoll(state, indices)
		if !system.IsValidation(err) {
			t.Fatalf("accepted indices %v", indices)
		}
	}
}
