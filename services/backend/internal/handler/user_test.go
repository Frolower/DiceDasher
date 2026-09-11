package handler

import (
	"backend/internal/auth"
	"context"
	"diceDasher/pkg/httputil"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type creatorFunc func(context.Context, auth.CreateInput) (auth.Created, error)

func (f creatorFunc) Register(ctx context.Context, input auth.CreateInput) (auth.Created, error) {
	return f(ctx, input)
}

func TestRegistrationHTTPContract(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		err        error
		status     int
		called     bool
	}{
		{"created", `{"username":"Alice","email":"alice@example.com","password":"StrongPass1"}`, nil, 201, true},
		{"invalid input", `{}`, auth.InvalidInputError{Err: errors.New("invalid fields")}, 400, true},
		{"duplicate", `{}`, auth.ErrAlreadyExists, 409, true},
		{"storage", `{}`, errors.New("private database details"), 500, true},
		{"empty", ``, nil, 400, false},
		{"malformed", `{`, nil, 400, false},
		{"unknown field", `{"extra":true}`, nil, 400, false},
		{"wrong type", `{"password":123}`, nil, 400, false},
		{"multiple values", `{} {}`, nil, 400, false},
		{"oversized", `{"username":"` + strings.Repeat("x", 17000) + `"}`, nil, 413, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			h := New(creatorFunc(func(ctx context.Context, input auth.CreateInput) (auth.Created, error) {
				called = true
				if tc.name == "created" && (input.Username != "Alice" || input.Email != "alice@example.com" || input.Password != "StrongPass1") {
					t.Fatal("lost request fields")
				}
				return auth.Created{}, tc.err
			}))
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tc.body)))
			if w.Code != tc.status || called != tc.called {
				t.Fatalf("status %d called %v body %s", w.Code, called, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "private") {
				t.Fatal("leaked internal error")
			}
			if tc.status == 201 && (w.Body.String() != "{\"status\":\"created\"}\n" || w.Header().Get("Content-Type") != "application/json") {
				t.Fatalf("unexpected success: %s", w.Body.String())
			}
			if len(w.Result().Cookies()) != 0 {
				t.Fatal("registration created cookie")
			}
		})
	}
}
