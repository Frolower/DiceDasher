package handler

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/service"
	"diceDasher/services/resolve/internal/system"
	"diceDasher/services/resolve/internal/system/generic"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type serviceFunc func(context.Context, service.Command) (service.Result, error)

func (f serviceFunc) Resolve(ctx context.Context, cmd service.Command) (service.Result, error) {
	return f(ctx, cmd)
}
func TestErrorMapping(t *testing.T) {
	for _, tc := range []struct {
		err     error
		status  int
		message string
	}{
		{system.Invalid(system.BadRequest, errors.New("invalid action")), 400, "invalid action"},
		{system.Invalid(system.Validation, errors.New("bad pool")), 422, "bad pool"},
		{system.ErrUnknownSystem, 404, "unknown system"},
		{repository.ErrNotFound, 404, "record not found"},
		{system.ErrInvalidTransition, 409, system.ErrInvalidTransition.Error()},
		{system.ErrLegacyContinuation, 409, system.ErrLegacyContinuation.Error()},
		{errors.New("postgres password=secret"), 500, "internal error"},
	} {
		h := New(serviceFunc(func(context.Context, service.Command) (service.Result, error) {
			return service.Result{}, fmt.Errorf("operation: %w", tc.err)
		}))
		w := httptest.NewRecorder()
		h.Resolve(w, httptest.NewRequest(http.MethodPost, "/resolve?system=tes", strings.NewReader(`{}`)))
		if w.Code != tc.status || strings.TrimSpace(w.Body.String()) != tc.message {
			t.Fatalf("response %d %s", w.Code, w.Body.String())
		}
	}
}
func TestHTTPCommandAndEnvelope(t *testing.T) {
	id := uuid.New()
	h := New(serviceFunc(func(_ context.Context, cmd service.Command) (service.Result, error) {
		if cmd.System != "tes" || cmd.Action != "push" || string(cmd.Payload) != `{"record_id":"example"}` {
			t.Fatalf("wrong command: %+v", cmd)
		}
		return service.Result{RecordID: id, Payload: json.RawMessage(`{"success":true}`)}, nil
	}))
	w := httptest.NewRecorder()
	h.Resolve(w, httptest.NewRequest(http.MethodPost, "/resolve?system=tes&action=push", strings.NewReader(`{"record_id":"example"}`)))
	var result service.Result
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || result.RecordID != id || string(result.Payload) != `{"success":true}` {
		t.Fatalf("response %d %s", w.Code, w.Body.String())
	}
}
func TestMalformedJSONStopsBeforeService(t *testing.T) {
	h := New(serviceFunc(func(context.Context, service.Command) (service.Result, error) {
		t.Fatal("service called")
		return service.Result{}, nil
	}))
	w := httptest.NewRecorder()
	h.Resolve(w, httptest.NewRequest(http.MethodPost, "/resolve?system=tes", strings.NewReader(`{`)))
	if w.Code != 400 {
		t.Fatal(w.Code)
	}
}

type failingHistory struct{}

func (failingHistory) InsertRollHistory(context.Context, repository.RollHistory) (uuid.UUID, error) {
	return uuid.Nil, errors.New("database connection secret")
}
func TestStorageFailureReturns500(t *testing.T) {
	svc := service.New(map[string]system.Resolver{"generic": generic.Resolver{Dice: dice.NewGenerator(func(int) int { return 0 })}}, failingHistory{})
	h := New(svc)
	w := httptest.NewRecorder()
	h.Resolve(w, httptest.NewRequest(http.MethodPost, "/resolve?system=generic", strings.NewReader(`{"number":1,"size":6}`)))
	if w.Code != 500 || strings.TrimSpace(w.Body.String()) != "internal error" {
		t.Fatalf("response %d %s", w.Code, w.Body.String())
	}
}
