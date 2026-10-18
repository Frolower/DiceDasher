package handler

import (
	"backend/internal/auth"
	"context"
	"diceDasher/pkg/httputil"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type authenticateFunc struct {
	Authenticator
	authenticate func(context.Context, string) (auth.AccessClaims, error)
}

func (f authenticateFunc) Authenticate(ctx context.Context, token string) (auth.AccessClaims, error) {
	return f.authenticate(ctx, token)
}

func TestAuthMiddleware(t *testing.T) {
	for _, tc := range []struct {
		name    string
		headers []string
		err     error
		status  int
		called  bool
	}{
		{name: "valid", headers: []string{"Bearer access-secret"}, status: 204, called: true},
		{name: "case insensitive scheme", headers: []string{"bEaReR access-secret"}, status: 204, called: true},
		{name: "missing", status: 401},
		{name: "empty", headers: []string{""}, status: 401},
		{name: "missing token", headers: []string{"Bearer"}, status: 401},
		{name: "wrong scheme", headers: []string{"Basic access-secret"}, status: 401},
		{name: "extra fields", headers: []string{"Bearer access-secret extra"}, status: 401},
		{name: "duplicate headers", headers: []string{"Bearer access-secret", "Bearer other"}, status: 401},
		{name: "invalid JWT", headers: []string{"Bearer access-secret"}, err: fmt.Errorf("verify: %w", auth.ErrInvalidToken), status: 401, called: true},
		{name: "invalid session", headers: []string{"Bearer access-secret"}, err: fmt.Errorf("lookup: %w", auth.ErrSessionNotFound), status: 401, called: true},
		{name: "database failure", headers: []string{"Bearer access-secret"}, err: errors.New("private database details"), status: 500, called: true},
		{name: "unconfigured", headers: []string{"Bearer access-secret"}, err: auth.ErrAuthNotConfigured, status: 500, called: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			claims := auth.AccessClaims{UserID: uuid.New(), SessionID: uuid.New(), ExpiresAt: time.Now().Add(time.Minute)}
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			for _, header := range tc.headers {
				req.Header.Add("Authorization", header)
			}
			// Cancellation and existing values must survive adding authentication claims.
			type originalKey struct{}
			ctx, cancel := context.WithCancel(context.WithValue(req.Context(), originalKey{}, "existing"))
			defer cancel()
			req = req.WithContext(ctx)
			calls, nextCalls := 0, 0
			h := New(authenticateFunc{authenticate: func(received context.Context, token string) (auth.AccessClaims, error) {
				calls++
				if received != ctx || token != "access-secret" {
					t.Fatal("wrong authentication arguments")
				}
				return claims, tc.err
			}})
			next := func(w http.ResponseWriter, r *http.Request) {
				nextCalls++
				got, ok := AccessClaimsFromContext(r.Context())
				if !ok || got != claims {
					t.Fatalf("missing or wrong claims: %v", got)
				}
				if r.Context().Value(originalKey{}) != "existing" {
					t.Fatal("lost original context")
				}
				cancel()
				if r.Context().Err() != context.Canceled {
					t.Fatal("lost cancellation")
				}
				w.WriteHeader(http.StatusNoContent)
			}
			w := httptest.NewRecorder()
			h.RequreAuth(next)(w, req)
			if w.Code != tc.status || (calls == 1) != tc.called || calls > 1 {
				t.Fatalf("status=%d calls=%d body=%s", w.Code, calls, w.Body.String())
			}
			wantNext := 0
			if tc.status == 204 {
				wantNext = 1
			}
			if nextCalls != wantNext {
				t.Fatalf("next called %d times", nextCalls)
			}
			if _, ok := AccessClaimsFromContext(ctx); ok {
				t.Fatal("claims leaked into original context")
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Pragma") != "no-cache" {
				t.Fatal("cacheable response")
			}
			challenge := w.Header().Get("WWW-Authenticate")
			if tc.status == 401 && challenge != `Bearer realm="backend"` {
				t.Fatalf("wrong challenge: %q", challenge)
			}
			if tc.status != 401 && challenge != "" {
				t.Fatal("unexpected challenge")
			}
			if strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "access-secret") {
				t.Fatal("leaked secrets")
			}
		})
	}
}
func TestAccessClaimsFromEmptyContext(t *testing.T) {
	claims, ok := AccessClaimsFromContext(context.Background())
	if ok || claims != (auth.AccessClaims{}) {
		t.Fatal("unexpected claims")
	}
}

type middlewareSessionStore struct {
	auth.SessionStore
	session auth.Session
}

func (s *middlewareSessionStore) FindSessionByID(_ context.Context, id uuid.UUID) (auth.Session, error) {
	if id != s.session.ID {
		return auth.Session{}, auth.ErrSessionNotFound
	}
	return s.session, nil
}

func TestAuthMiddlewareWithJWTAndSessionRevocation(t *testing.T) {
	manager, err := auth.NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	store := &middlewareSessionStore{session: auth.Session{
		ID: uuid.New(), UserID: uuid.New(), ExpiresAt: now.Add(time.Hour),
	}}
	raw, err := manager.IssueAccess(store.session.UserID, store.session.ID, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	h := New(auth.NewServiceWithDependencies(nil, store, nil, manager))
	calls := 0
	protected := h.RequreAuth(func(w http.ResponseWriter, r *http.Request) {
		calls++
		claims, ok := AccessClaimsFromContext(r.Context())
		if !ok || claims.UserID != store.session.UserID || claims.SessionID != store.session.ID {
			t.Fatal("wrong authenticated identity")
		}
		w.WriteHeader(http.StatusNoContent)
	})
	request := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+raw)
		w := httptest.NewRecorder()
		protected(w, req)
		return w
	}
	if w := request(); w.Code != http.StatusNoContent || calls != 1 {
		t.Fatalf("active session rejected: %d %s", w.Code, w.Body.String())
	}
	store.session.RevokedAt = &now
	if w := request(); w.Code != http.StatusUnauthorized || calls != 1 {
		t.Fatal("revoked session accepted an otherwise valid JWT")
	}
}
func TestProtectedRoutesRequireAuth(t *testing.T) {
	h := New(authenticateFunc{authenticate: func(context.Context, string) (auth.AccessClaims, error) {
		return auth.AccessClaims{UserID: uuid.New(), SessionID: uuid.New()}, nil
	}})
	router := httputil.NewRouter()
	h.RegisterRouters(router)
	for _, route := range []struct{ method, path string }{
		{http.MethodGet, "/me"},
		{http.MethodPatch, "/user"},
		{http.MethodDelete, "/user"},
		{http.MethodPut, "/password"},
	} {
		t.Run(route.method+route.path, func(t *testing.T) {
			for _, authorized := range []bool{false, true} {
				req := httptest.NewRequest(route.method, route.path, nil)
				want := http.StatusUnauthorized
				if authorized {
					req.Header.Set("Authorization", "Bearer valid")
					want = http.StatusNotImplemented
				}
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				if w.Code != want {
					t.Fatalf("authorized=%v status=%d want=%d", authorized, w.Code, want)
				}
			}
		})
	}
}
