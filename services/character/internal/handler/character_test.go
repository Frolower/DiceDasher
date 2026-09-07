package handler

import (
	"context"
	"diceDasher/pkg/httputil"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"diceDasher/services/character/internal/service"
	"diceDasher/services/character/internal/system"
	"diceDasher/services/character/internal/system/tes"
)

func TestCharacterValidationResponse(t *testing.T) {
	h := New(service.New(nil, map[string]system.Character{"tes": tes.Character{}}))
	req := httptest.NewRequest(http.MethodPost, "/character?system=tes", strings.NewReader(`{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":{"name":""}}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.postUserCreatedCharacterHandler(w, req)
	if w.Code != 422 {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Fatal("expected JSON")
	}
	var response system.ValidationError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Violations) != 1 || response.Violations[0].Field != "character.name" || response.Violations[0].Code != "required" {
		t.Fatalf("unexpected violations: %+v", response)
	}
}

type creatorFunc func(context.Context, string, json.RawMessage) (uuid.UUID, error)

func (f creatorFunc) Create(ctx context.Context, name string, raw json.RawMessage) (uuid.UUID, error) {
	return f(ctx, name, raw)
}

func TestCharacterHTTPContract(t *testing.T) {
	id := uuid.New()
	for _, tt := range []struct {
		name, url, body string
		err             error
		status          int
		called          bool
		message         string
	}{
		{name: "created", url: "/character?system=tes", body: `{}`, status: 201, called: true, message: id.String()},
		{name: "missing system", url: "/character", body: `{}`, status: 400, message: "missing query param"},
		{name: "malformed JSON", url: "/character?system=tes", body: `{`, status: 400},
		{name: "unknown system", url: "/character?system=other", body: `{}`, err: service.ErrUnknownSystem, status: 404, called: true, message: "unknown system"},
		{name: "input", url: "/character?system=tes", body: `{}`, err: fmt.Errorf("create: %w", &system.InputError{Err: errors.New("user_id is required")}), status: 400, called: true, message: "user_id is required"},
		{name: "storage", url: "/character?system=tes", body: `{}`, err: &service.SaveError{Err: errors.New("private connection details")}, status: 500, called: true, message: "failed to save character"},
		{name: "internal", url: "/character?system=tes", body: `{}`, err: errors.New("private connection details"), status: 500, called: true, message: "internal error"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			h := New(creatorFunc(func(ctx context.Context, name string, raw json.RawMessage) (uuid.UUID, error) {
				called = true
				return id, tt.err
			}))
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.body)))
			if w.Code != tt.status || called != tt.called || !strings.Contains(w.Body.String(), tt.message) {
				t.Fatalf("status %d, called %v, body %s", w.Code, called, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "private") {
				t.Fatal("internal error leaked")
			}
			if tt.status == 201 {
				var response createCharacterResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.ID != id {
					t.Fatalf("invalid success response: %s", w.Body.String())
				}
			}
		})
	}
}
