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
	"time"
)

type refreshMemoryStore struct {
	auth.SessionStore
	session auth.Session
}

func (s *refreshMemoryStore) FindSessionByRefreshHash(_ context.Context, hash string) (auth.Session, error) {
	if hash != s.session.RefreshTokenHash {
		return auth.Session{}, auth.ErrSessionNotFound
	}
	return s.session, nil
}

func (s *refreshMemoryStore) RotateRefreshToken(_ context.Context, id uuid.UUID, oldHash, newHash string, now time.Time) error {
	if id != s.session.ID || oldHash != s.session.RefreshTokenHash || s.session.RevokedAt != nil || !s.session.ExpiresAt.After(now) {
		return auth.ErrSessionNotFound
	}
	s.session.RefreshTokenHash = newHash
	return nil
}

func TestRefreshHTTPWithAuthService(t *testing.T) {
	manager, err := auth.NewJWTManager([]byte(strings.Repeat("k", 32)), "issuer", "audience")
	if err != nil {
		t.Fatal(err)
	}
	store := &refreshMemoryStore{session: auth.Session{
		ID: uuid.New(), UserID: uuid.New(), RefreshTokenHash: manager.HashRefresh("original-refresh"),
		ExpiresAt: time.Now().UTC().Truncate(time.Second).Add(time.Hour),
	}}
	service := auth.NewServiceWithDependencies(nil, store, nil, manager)
	h := NewWithOptions(service, Options{CookieSecure: true})
	router := httputil.NewRouter()
	h.RegisterRouters(router)
	refresh := func(token string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "__Host-refresh_token", Value: token})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	first := refresh("original-refresh")
	if first.Code != http.StatusOK || len(first.Result().Cookies()) != 1 {
		t.Fatalf("refresh failed: %d %s", first.Code, first.Body.String())
	}
	rotated := first.Result().Cookies()[0].Value
	if store.session.RefreshTokenHash != manager.HashRefresh(rotated) || rotated == "original-refresh" {
		t.Fatal("new cookie does not match stored hash")
	}
	var response tokenResponse
	if err := json.Unmarshal(first.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	claims, err := manager.VerifyAccess(response.AccessToken)
	if err != nil || claims.UserID != store.session.UserID || claims.SessionID != store.session.ID {
		t.Fatalf("wrong access JWT: %v %v", claims, err)
	}
	if reused := refresh("original-refresh"); reused.Code != http.StatusUnauthorized || len(reused.Result().Cookies()) != 0 {
		t.Fatal("old refresh token remained usable")
	}
	if next := refresh(rotated); next.Code != http.StatusOK {
		t.Fatalf("new refresh token rejected: %d %s", next.Code, next.Body.String())
	}
}

type refreshFunc struct {
	Authenticator
	refresh func(context.Context, string) (auth.Tokens, error)
}

func (f refreshFunc) Refresh(ctx context.Context, token string) (auth.Tokens, error) {
	return f.refresh(ctx, token)
}
func TestRefreshHTTP(t *testing.T) {
	for _, tc := range []struct {
		name, origin, fetchSite             string
		secure, cookie, wrongCookie, called bool
		err                                 error
		status                              int
	}{
		{name: "secure success without body", secure: true, cookie: true, called: true, status: 200},
		{name: "local success", cookie: true, called: true, status: 200},
		{name: "missing cookie", secure: true, status: 401},
		{name: "wrong cookie name", secure: true, cookie: true, wrongCookie: true, status: 401},
		{name: "invalid session", cookie: true, called: true, err: fmt.Errorf("lookup: %w", auth.ErrSessionNotFound), status: 401},
		{name: "database failure", cookie: true, called: true, err: errors.New("private database details"), status: 500},
		{name: "unconfigured", cookie: true, called: true, err: auth.ErrAuthNotConfigured, status: 500},
		{name: "allowed origin", origin: "https://app.example", cookie: true, called: true, status: 200},
		{name: "foreign origin", origin: "https://evil.example", cookie: true, status: 403},
		{name: "opaque origin", origin: "null", cookie: true, status: 403},
		{name: "cross site without origin", fetchSite: "cross-site", cookie: true, status: 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			now := time.Now().UTC().Truncate(time.Second)
			tokens := auth.Tokens{AccessToken: "new-access", RefreshToken: "secret-new-refresh", AccessExpiresAt: now.Add(time.Minute), RefreshExpiresAt: now.Add(time.Hour)}
			h := NewWithOptions(refreshFunc{refresh: func(ctx context.Context, token string) (auth.Tokens, error) {
				calls++
				if token != "secret-old-refresh" {
					t.Fatalf("wrong token: %q", token)
				}
				return tokens, tc.err
			}}, Options{CookieSecure: tc.secure, AllowedOrigin: "https://app.example"})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
			req.Header.Set("Origin", tc.origin)
			req.Header.Set("Sec-Fetch-Site", tc.fetchSite)
			name := "refresh_token"
			if tc.secure {
				name = "__Host-refresh_token"
			}
			if tc.cookie {
				cookieName := name
				if tc.wrongCookie {
					cookieName = "refresh_token"
				}
				req.AddCookie(&http.Cookie{Name: cookieName, Value: "secret-old-refresh"})
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.status || (calls == 1) != tc.called || calls > 1 {
				t.Fatalf("status=%d calls=%d body=%s", w.Code, calls, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Pragma") != "no-cache" {
				t.Fatal("cacheable response")
			}
			if strings.Contains(w.Body.String(), "secret-") || strings.Contains(w.Body.String(), "private") {
				t.Fatal("leaked secrets")
			}
			cookies := w.Result().Cookies()
			if tc.status != 200 {
				if len(cookies) != 0 || strings.Contains(w.Body.String(), "new-access") {
					t.Fatal("tokens returned on failure")
				}
				return
			}
			if len(cookies) != 1 {
				t.Fatal("missing rotated cookie")
			}
			c := cookies[0]
			if c.Name != name || c.Value != tokens.RefreshToken || c.Path != "/" || c.Domain != "" || !c.HttpOnly || c.Secure != tc.secure || c.SameSite != http.SameSiteStrictMode || c.MaxAge <= 0 || !c.Expires.Equal(tokens.RefreshExpiresAt) {
				t.Fatalf("wrong cookie: %v", c)
			}
			var response tokenResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.AccessToken != tokens.AccessToken || response.TokenType != "Bearer" || !response.ExpiresAt.Equal(tokens.AccessExpiresAt) {
				t.Fatalf("wrong response: %s", w.Body.String())
			}
		})
	}
}
