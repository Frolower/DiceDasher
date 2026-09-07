package service

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
	"diceDasher/services/resolve/internal/system/tes"
	"diceDasher/services/resolve/internal/system/vtmv5"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"testing"
)

type memoryHistory struct {
	records map[uuid.UUID]repository.RollHistory
	failure error
	writes  int
}

func (m *memoryHistory) InsertRollHistory(ctx context.Context, rec repository.RollHistory) (uuid.UUID, error) {
	m.writes++
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}
	if m.failure != nil {
		return uuid.Nil, m.failure
	}
	rec.ID = uuid.New()
	m.records[rec.ID] = rec
	return rec.ID, nil
}
func (m *memoryHistory) GetRollHistoryByID(ctx context.Context, id uuid.UUID) (repository.RollHistory, error) {
	if err := ctx.Err(); err != nil {
		return repository.RollHistory{}, err
	}
	rec, ok := m.records[id]
	if !ok {
		return rec, repository.ErrNotFound
	}
	return rec, nil
}

type resolverFunc func(context.Context, string, json.RawMessage) (any, error)

func (f resolverFunc) Resolve(ctx context.Context, action string, raw json.RawMessage) (any, error) {
	return f(ctx, action, raw)
}

func TestSavedContinuations(t *testing.T) {
	for _, name := range []string{"tes", "vtmv5"} {
		t.Run(name, func(t *testing.T) {
			history := &memoryHistory{records: map[uuid.UUID]repository.RollHistory{}}
			generator := dice.NewGenerator(func(n int) int { return n - 1 })
			resolvers := map[string]system.Resolver{"tes": tes.New(history, generator), "vtmv5": vtmv5.New(history, generator)}
			svc := New(resolvers, history)
			// The service owns a registry snapshot, independent of caller mutations.
			delete(resolvers, name)
			payload := json.RawMessage(`{"attr":2,"gear":1,"target":2}`)
			action := "push"
			if name == "vtmv5" {
				payload = json.RawMessage(`{"attribute":3,"hunger":1,"target":2}`)
				action = "reroll"
			}
			requestID := uuid.New()
			result, err := svc.Resolve(context.Background(), Command{System: name, Payload: payload, RequestID: requestID})
			if err != nil {
				t.Fatal(err)
			}
			root := result.RecordID
			rec := history.records[root]
			if root == uuid.Nil || rec.ActionType != "roll" || rec.RequestID != requestID || string(rec.ResponsePayload) != string(result.Payload) || len(rec.StatePayload) == 0 {
				t.Fatalf("bad history: %+v", rec)
			}
			for i := 0; i < 2; i++ {
				parent := result.RecordID
				raw, _ := json.Marshal(map[string]any{"record_id": parent, "reroll_index": []int{0}})
				result, err = svc.Resolve(context.Background(), Command{System: name, Action: action, Payload: raw})
				if err != nil {
					t.Fatal(err)
				}
				rec = history.records[result.RecordID]
				var state struct {
					OriginalID uuid.UUID `json:"original_id"`
					ParentID   uuid.UUID `json:"parent_id"`
					Target     int       `json:"target"`
				}
				if err := json.Unmarshal(rec.StatePayload, &state); err != nil {
					t.Fatal(err)
				}
				if state.OriginalID != root || state.ParentID != parent || state.Target != 2 || rec.RequestID == uuid.Nil {
					t.Fatalf("lost continuation: %+v", state)
				}
			}
		})
	}
}

func TestFailuresDoNotReturnSuccessOrReroll(t *testing.T) {
	storageError := errors.New("database unavailable")
	calls := 0
	resolver := resolverFunc(func(context.Context, string, json.RawMessage) (any, error) {
		calls++
		return map[string]int{"roll": 6}, nil
	})
	history := &memoryHistory{failure: storageError}
	svc := New(map[string]system.Resolver{"test": resolver}, history)
	result, err := svc.Resolve(context.Background(), Command{System: "test"})
	if !errors.Is(err, storageError) || result.RecordID != uuid.Nil || result.Payload != nil || calls != 1 || history.writes != 1 {
		t.Fatalf("storage failure: %+v %v calls=%d writes=%d", result, err, calls, history.writes)
	}
}

func TestFailedExecutionIsNotSaved(t *testing.T) {
	invalid := system.Invalid(system.Validation, errors.New("bad pool"))
	for _, response := range []any{nil, make(chan int), badState{}} {
		history := &memoryHistory{records: map[uuid.UUID]repository.RollHistory{}}
		resolver := resolverFunc(func(context.Context, string, json.RawMessage) (any, error) {
			if response == nil {
				return nil, invalid
			}
			return response, nil
		})
		svc := New(map[string]system.Resolver{"test": resolver}, history)
		if result, err := svc.Resolve(context.Background(), Command{System: "test"}); err == nil || result.RecordID != uuid.Nil || history.writes != 0 {
			t.Fatalf("saved failed execution: %+v %v", result, err)
		}
	}
}

type badState struct{}

func (badState) HistoryState() any { return make(chan int) }

func TestDispatchErrors(t *testing.T) {
	svc := New(nil, &memoryHistory{})
	if _, err := svc.Resolve(context.Background(), Command{System: "unknown"}); !errors.Is(err, system.ErrUnknownSystem) {
		t.Fatal(err)
	}
	if _, err := svc.Resolve(context.Background(), Command{}); err == nil {
		t.Fatal("missing system accepted")
	}
}
