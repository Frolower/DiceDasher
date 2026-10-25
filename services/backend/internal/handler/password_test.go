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

type passwordUpdateFunc struct {
	User
	update func(context.Context, user.PasswordUpdateRecord) error
}

func (f passwordUpdateFunc) UpdatePassword(ctx context.Context, rec user.PasswordUpdateRecord) error {
	return f.update(ctx, rec)
}

func TestPasswordUpdateHTTPContract(t *testing.T) {
	const validBody = `{"old_password":"OldPass1!","new_password":"NewPass2!"}`
	for _, tc := range []struct {
		name, body string
		err        error
		status     int
		called     bool
	}{
		{"success", validBody, nil, 200, true},
		{"invalid input", validBody, user.InvalidInputError{Err: errors.New("invalid password")}, 400, true},
		{"wrong old password", validBody, fmt.Errorf("verify: %w", user.ErrOldPasswordDoesntMatch), 400, true},
		{"missing user", validBody, fmt.Errorf("lookup: %w", user.ErrUserNotFound), 404, true},
		{"database failure", validBody, errors.New("private database details"), 500, true},
		{"empty body", "", nil, 400, false},
		{"malformed JSON", "{", nil, 400, false},
		{"unknown field", `{"password":"secret"}`, nil, 400, false},
		{"ID in body", `{"id":"` + uuid.NewString() + `"}`, nil, 400, false},
		{"wrong type", `{"new_password":123}`, nil, 400, false},
		{"multiple values", validBody + `{}`, nil, 400, false},
		{"oversized", `{"new_password":"` + strings.Repeat("x", 17000) + `"}`, nil, 413, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			calls := 0
			h := New(authenticateFunc{authenticate: func(_ context.Context, token string) (auth.AccessClaims, error) {
				if token != "secret-token" {
					t.Fatal("wrong token")
				}
				return auth.AccessClaims{UserID: id, SessionID: uuid.New()}, nil
			}}, passwordUpdateFunc{update: func(ctx context.Context, rec user.PasswordUpdateRecord) error {
				calls++
				claims, ok := AccessClaimsFromContext(ctx)
				if !ok || claims.UserID != id || rec != (user.PasswordUpdateRecord{Id: id, OldPassword: "OldPass1!", NewPassword: "NewPass2!"}) {
					t.Fatal("lost password fields or authenticated identity")
				}
				return tc.err
			}})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPut, "/password?id="+uuid.NewString(), strings.NewReader(tc.body))
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
			for _, secret := range []string{"private", "secret-token", "OldPass1!", "NewPass2!"} {
				if strings.Contains(w.Body.String(), secret) {
					t.Fatal("response leaked private data")
				}
			}
			if tc.status == http.StatusOK && (w.Body.String() != "\"success\"\n" || w.Header().Get("Content-Type") != "application/json") {
				t.Fatalf("unexpected success response: %s", w.Body.String())
			}
		})
	}
}

func TestPasswordUpdateRejectsUnauthenticatedRequests(t *testing.T) {
	for _, tc := range []struct {
		name, token string
		err         error
	}{
		{name: "missing token"},
		{name: "invalid token", token: "invalid", err: auth.ErrInvalidToken},
		{name: "revoked session", token: "revoked", err: auth.ErrSessionNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := New(authenticateFunc{authenticate: func(context.Context, string) (auth.AccessClaims, error) {
				return auth.AccessClaims{}, tc.err
			}}, passwordUpdateFunc{update: func(context.Context, user.PasswordUpdateRecord) error {
				t.Fatal("unauthenticated request reached update")
				return nil
			}})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPut, "/password", strings.NewReader(`{"old_password":"OldPass1!","new_password":"NewPass2!"}`))
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status=%d", w.Code)
			}
		})
	}
}
