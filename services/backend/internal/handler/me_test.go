package handler

import (
	"backend/internal/auth"
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
)

type profileFunc struct {
	Authenticator
	authenticate func(context.Context, string) (auth.AccessClaims, error)
	profile      func(context.Context, uuid.UUID) (auth.UserProfile, error)
}

func (f profileFunc) Authenticate(ctx context.Context, token string) (auth.AccessClaims, error) {
	return f.authenticate(ctx, token)
}
func (f profileFunc) GetProfile(ctx context.Context, id uuid.UUID) (auth.UserProfile, error) {
	return f.profile(ctx, id)
}

func TestMeHTTP(t *testing.T) {
	for _, tc := range []struct {
		name                            string
		authorized                      bool
		authErr, profileErr             error
		status, authCalls, profileCalls int
	}{
		{name: "success", authorized: true, status: 200, authCalls: 1, profileCalls: 1},
		{name: "missing token", status: 401},
		{name: "invalid token", authorized: true, authErr: auth.ErrInvalidToken, status: 401, authCalls: 1},
		{name: "deleted user", authorized: true, profileErr: fmt.Errorf("lookup: %w", auth.ErrUserNotFound), status: 401, authCalls: 1, profileCalls: 1},
		{name: "database failure", authorized: true, profileErr: errors.New("private database details"), status: 500, authCalls: 1, profileCalls: 1},
		{name: "unconfigured", authorized: true, profileErr: auth.ErrAuthNotConfigured, status: 500, authCalls: 1, profileCalls: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			profile := auth.UserProfile{ID: uuid.New(), Username: "Alice", Email: "alice@example.com"}
			authCalls, profileCalls := 0, 0
			h := New(profileFunc{
				authenticate: func(ctx context.Context, token string) (auth.AccessClaims, error) {
					authCalls++
					if token != "secret-token" {
						t.Fatalf("wrong token: %s", token)
					}
					return auth.AccessClaims{UserID: profile.ID, SessionID: uuid.New()}, tc.authErr
				},
				profile: func(ctx context.Context, id uuid.UUID) (auth.UserProfile, error) {
					profileCalls++
					if id != profile.ID {
						t.Fatal("handler did not use authenticated user ID")
					}
					if claims, ok := AccessClaimsFromContext(ctx); !ok || claims.UserID != id {
						t.Fatal("missing identity in profile lookup context")
					}
					return profile, tc.profileErr
				},
			})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			// Untrusted IDs must not override the identity from middleware.
			req := httptest.NewRequest(http.MethodGet, "/me?id="+uuid.NewString(), nil)
			if tc.authorized {
				req.Header.Set("Authorization", "Bearer secret-token")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.status || authCalls != tc.authCalls || profileCalls != tc.profileCalls {
				t.Fatalf("status=%d authCalls=%d profileCalls=%d body=%s", w.Code, authCalls, profileCalls, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Pragma") != "no-cache" {
				t.Fatal("cacheable profile")
			}
			if len(w.Result().Cookies()) != 0 {
				t.Fatal("profile endpoint changed cookies")
			}
			if strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "secret-token") {
				t.Fatal("leaked secrets")
			}
			if tc.status == 401 && w.Header().Get("WWW-Authenticate") == "" {
				t.Fatal("missing authentication challenge")
			}
			if tc.status != 200 {
				if strings.Contains(w.Body.String(), profile.Email) {
					t.Fatal("profile returned on failure")
				}
				return
			}
			if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
				t.Fatal("wrong content type")
			}
			var body map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if len(body) != 3 || body["id"] != profile.ID.String() || body["username"] != profile.Username || body["email"] != profile.Email {
				t.Fatalf("wrong public profile: %v", body)
			}
		})
	}
}
func TestMeWithoutAuthenticationContext(t *testing.T) {
	h := New(profileFunc{profile: func(context.Context, uuid.UUID) (auth.UserProfile, error) {
		t.Fatal("lookup without authenticated identity")
		return auth.UserProfile{}, nil
	}})
	w := httptest.NewRecorder()
	h.HandleMe(w, httptest.NewRequest(http.MethodGet, "/me", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w.Code)
	}
}
