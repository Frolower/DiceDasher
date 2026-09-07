package service

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"diceDasher/services/character/internal/system"
	"diceDasher/services/character/internal/system/tes"
	"github.com/google/uuid"
)

type repositoryFunc func(context.Context, Record) (uuid.UUID, error)

func (f repositoryFunc) InsertPlayerCreatedCharacter(ctx context.Context, rec Record) (uuid.UUID, error) {
	return f(ctx, rec)
}

type strategyFunc func(context.Context, json.RawMessage) (system.CreatedCharacter, error)

func (f strategyFunc) CreateCharacter(ctx context.Context, raw json.RawMessage) (system.CreatedCharacter, error) {
	return f(ctx, raw)
}

func TestCreateSelectsStrategyAndPersistsRecord(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	raw := json.RawMessage(`{"name":"Mira"}`)
	created := system.CreatedCharacter{UserID: uuid.New(), CharacterType: "pc", Name: "Mira", Data: raw}
	wantID := uuid.New()
	calls := 0
	systems := map[string]system.Character{"chosen": strategyFunc(func(got context.Context, body json.RawMessage) (system.CreatedCharacter, error) {
		if got != ctx || string(body) != string(raw) {
			t.Fatal("lost context or payload")
		}
		return created, nil
	})}
	s := New(repositoryFunc(func(got context.Context, rec Record) (uuid.UUID, error) {
		calls++
		want := Record{UserID: created.UserID, SystemName: "chosen", CharacterType: "pc", Name: "Mira", Data: raw}
		if got != ctx || !reflect.DeepEqual(rec, want) {
			t.Fatalf("unexpected record: %+v", rec)
		}
		return wantID, nil
	}), systems)
	delete(systems, "chosen") // The constructor owns its registry snapshot.
	id, err := s.Create(ctx, "chosen", raw)
	if err != nil || id != wantID || calls != 1 {
		t.Fatalf("id %s, calls %d, err %v", id, calls, err)
	}
}

func TestCreateDoesNotSaveRejectedInput(t *testing.T) {
	s := New(repositoryFunc(func(context.Context, Record) (uuid.UUID, error) {
		t.Fatal("repository called for invalid input")
		return uuid.Nil, nil
	}), map[string]system.Character{"tes": tes.Character{}})
	for _, tt := range []struct{ name, body string }{
		{"unknown", `{}`}, {"tes", `{`}, {"tes", `{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":{"name":""}}`},
	} {
		id, err := s.Create(context.Background(), tt.name, json.RawMessage(tt.body))
		if err == nil || id != uuid.Nil {
			t.Fatalf("accepted rejected request: %s", tt.body)
		}
		if tt.name == "unknown" && !errors.Is(err, ErrUnknownSystem) {
			t.Fatal(err)
		}
	}
}

func TestCreatePreservesFailures(t *testing.T) {
	validation := &system.ValidationError{Violations: []system.Violation{{Field: "character.name", Code: "required"}}}
	s := New(nil, map[string]system.Character{"test": strategyFunc(func(context.Context, json.RawMessage) (system.CreatedCharacter, error) {
		return system.CreatedCharacter{}, validation
	})})
	_, err := s.Create(context.Background(), "test", nil)
	var got *system.ValidationError
	if !errors.As(err, &got) || got != validation {
		t.Fatalf("lost domain error: %v", err)
	}
	cause := errors.New("database unavailable")
	s = New(repositoryFunc(func(context.Context, Record) (uuid.UUID, error) { return uuid.Nil, cause }), map[string]system.Character{"tes": tes.Character{}})
	id, err := s.Create(context.Background(), "tes", json.RawMessage(`{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":{"name":"Mira"}}`))
	var save *SaveError
	if id != uuid.Nil || !errors.As(err, &save) || !errors.Is(err, cause) {
		t.Fatalf("lost persistence error: %v", err)
	}
}
