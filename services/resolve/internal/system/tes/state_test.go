package tes

import (
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"testing"
)

func TestContinuationRoundTrip(t *testing.T) {
	root := uuid.New()
	rec := repository.RollHistory{ID: root, SystemName: "tes", ActionType: "roll",
		RequestPayload:  json.RawMessage(`{"target":2}`),
		ResponsePayload: json.RawMessage(`{"attribute_rolls":[1,6],"gear_rolls":[1,6]}`),
	}
	for i := 0; i < 3; i++ {
		state, err := loadState(rec)
		if err != nil {
			t.Fatal(err)
		}
		if state.Target != 2 || state.OriginalID != root || state.ParentID != rec.ID {
			t.Fatalf("lost state: %+v", state)
		}
		response, status, err := continueRoll(state)
		if err != nil || status != http.StatusOK {
			t.Fatalf("continuation: %d %v", status, err)
		}

		if !reflect.DeepEqual(state.GearRolls, []int{1, 6}) {
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
		rec = repository.RollHistory{ID: uuid.New(), SystemName: "tes", ActionType: "push", StatePayload: encoded,
			RequestPayload: json.RawMessage(`{"record_id":"` + rec.ID.String() + `"}`), ResponsePayload: payload}
	}
}

func TestRejectUnsupportedHistory(t *testing.T) {
	for _, tc := range []struct {
		system, action string
		expected       error
	}{
		{"generic", "roll", system.ErrInvalidTransition},
		{"tes", "check", system.ErrInvalidTransition},
		{"tes", "push", system.ErrLegacyContinuation},
	} {
		_, err := loadState(repository.RollHistory{SystemName: tc.system, ActionType: tc.action})
		if !errors.Is(err, tc.expected) {
			t.Fatalf("%s/%s: %v", tc.system, tc.action, err)
		}
	}
}

func TestInitialRollPersistsState(t *testing.T) {
	raw := json.RawMessage(`{"attr":2,"gear":1,"target":1}`)
	response, status, err := resolveRoll(raw)
	if err != nil || status != http.StatusOK {
		t.Fatalf("roll: %d %v", status, err)
	}
	encoded, err := json.Marshal(response.HistoryState())
	if err != nil {
		t.Fatal(err)
	}
	state, err := loadState(repository.RollHistory{ID: uuid.New(), SystemName: "tes", ActionType: "roll", StatePayload: encoded})
	if err != nil {
		t.Fatal(err)
	}
	if state.Target != 1 || !reflect.DeepEqual(state.AttributeRolls, response.AttributeRolls) || !reflect.DeepEqual(state.GearRolls, response.GearRolls) {
		t.Fatal("snapshot differs from result")
	}
}
