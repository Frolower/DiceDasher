package handler

import (
	"backend/internal/auth"
	"backend/internal/user"
	"context"
	"diceDasher/pkg/httputil"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type userUpdateFunc struct {
	User
	update func(context.Context, user.UserUpdateRecord) error
}

func (f userUpdateFunc) UpdateUser(ctx context.Context, rec user.UserUpdateRecord) error {
	return f.update(ctx, rec)
}

func TestUserUpdateHTTPContract(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		err        error
		status     int
		called     bool
	}{
		{"both fields", `{"username":"Alice","email":"alice@example.com"}`, nil, 200, true},
		{"username only", `{"username":"Alice"}`, nil, 200, true},
		{"email only", `{"email":"alice@example.com"}`, nil, 200, true},
		{"invalid input", `{}`, user.InvalidInputError{Err: errors.New("missing arguments to update")}, 400, true},
		{"duplicate", `{"username":"Alice"}`, fmt.Errorf("update: %w", user.ErrAlreadyExists), 409, true},
		{"not found", `{"username":"Alice"}`, fmt.Errorf("update: %w", user.ErrUserNotFound), 404, true},
		{"storage failure", `{"username":"Alice"}`, errors.New("private database details"), 500, true},
		{"empty body", ``, nil, 400, false},
		{"malformed JSON", `{`, nil, 400, false},
		{"unknown field", `{"extra":true}`, nil, 400, false},
		{"ID in body", `{"id":"` + uuid.NewString() + `","username":"Alice"}`, nil, 400, false},
		{"wrong type", `{"username":123}`, nil, 400, false},
		{"multiple values", `{} {}`, nil, 400, false},
		{"oversized", `{"username":"` + strings.Repeat("x", 17000) + `"}`, nil, 413, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			calls := 0
			h := New(authenticateFunc{authenticate: func(_ context.Context, token string) (auth.AccessClaims, error) {
				if token != "secret-token" {
					t.Fatal("wrong token")
				}
				return auth.AccessClaims{UserID: id, SessionID: uuid.New()}, nil
			}}, userUpdateFunc{update: func(ctx context.Context, rec user.UserUpdateRecord) error {
				calls++
				claims, ok := AccessClaimsFromContext(ctx)
				if !ok || claims.UserID != id || rec.ID != id {
					t.Fatal("update did not use authenticated identity")
				}
				wantUsername, wantEmail := "", ""
				if strings.Contains(tc.body, `"username":"Alice"`) {
					wantUsername = "Alice"
				}
				if strings.Contains(tc.body, `"email":"alice@example.com"`) {
					wantEmail = "alice@example.com"
				}
				if rec.Username != wantUsername || rec.Email != wantEmail {
					t.Fatalf("lost fields: %+v", rec)
				}
				return tc.err
			}})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPatch, "/user?id="+uuid.NewString(), strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer secret-token")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			wantCalls := 0
			if tc.called {
				wantCalls = 1
			}
			if w.Code != tc.status || calls != wantCalls {
				t.Fatalf("status=%d calls=%d body=%s", w.Code, calls, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "secret-token") {
				t.Fatal("leaked internal data")
			}
			if w.Header().Get("Cache-Control") != "no-store" || len(w.Result().Cookies()) != 0 {
				t.Fatal("unexpected cache or cookie behavior")
			}
			if tc.status == 200 && (w.Body.String() != "\"success\"\n" || w.Header().Get("Content-Type") != "application/json") {
				t.Fatalf("unexpected success: %s", w.Body.String())
			}
			if tc.name == "invalid input" && !strings.Contains(w.Body.String(), "missing arguments to update") {
				t.Fatal("lost validation detail")
			}
		})
	}
}

func TestUserUpdateRejectsUnauthenticatedRequests(t *testing.T) {
	for _, tc := range []struct {
		name, token string
		err         error
	}{
		{name: "missing token"},
		{name: "invalid token", token: "invalid", err: auth.ErrInvalidToken},
		{name: "revoked session", token: "revoked", err: auth.ErrSessionNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := New(authenticateFunc{authenticate: func(context.Context, string) (auth.AccessClaims, error) { return auth.AccessClaims{}, tc.err }}, userUpdateFunc{update: func(context.Context, user.UserUpdateRecord) error {
				t.Fatal("unauthenticated request reached update")
				return nil
			}})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPatch, "/user", strings.NewReader(`{"username":"Alice"}`))
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized || w.Header().Get("WWW-Authenticate") == "" {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestUserUpdateWithoutAuthenticationContext(t *testing.T) {
	h := New(nil, userUpdateFunc{update: func(context.Context, user.UserUpdateRecord) error {
		t.Fatal("update without authenticated identity")
		return nil
	}})
	w := httptest.NewRecorder()
	h.HandleUpdateUser(w, httptest.NewRequest(http.MethodPatch, "/user", strings.NewReader(`{"username":"Alice"}`)))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w.Code)
	}
}
