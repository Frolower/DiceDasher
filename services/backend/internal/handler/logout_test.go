package handler

import (
	"backend/internal/auth"
	"context"
	"diceDasher/pkg/httputil"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type logoutFunc struct {
	Authenticator
	logout func(context.Context, string) error
}

func (f logoutFunc) Logout(ctx context.Context, token string) error {
	return f.logout(ctx, token)
}

func TestLogoutHTTP(t *testing.T) {
	for _, tc := range []struct {
		name, origin, fetchSite, allowedOrigin string
		secure, cookie, called, cleared        bool
		err                                    error
		status                                 int
	}{
		{name: "secure success without body", secure: true, cookie: true, called: true, cleared: true, status: 204},
		{name: "local success", cookie: true, called: true, cleared: true, status: 204},
		{name: "no cookie", secure: true, status: 204},
		{name: "already revoked", secure: true, cookie: true, called: true, cleared: true, err: fmt.Errorf("revoke: %w", auth.ErrSessionNotFound), status: 204},
		{name: "database failure", secure: true, cookie: true, called: true, err: errors.New("private database details"), status: 500},
		{name: "unconfigured service", cookie: true, called: true, err: auth.ErrAuthNotConfigured, status: 500},
		{name: "allowed origin", origin: "https://app.example", allowedOrigin: "https://app.example", cookie: true, called: true, cleared: true, status: 204},
		{name: "foreign origin", origin: "https://evil.example", allowedOrigin: "https://app.example", cookie: true, status: 403},
		{name: "opaque origin", origin: "null", cookie: true, status: 403},
		{name: "origin without configuration", origin: "https://app.example", cookie: true, status: 403},
		{name: "cross site without origin", fetchSite: "cross-site", cookie: true, status: 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			h := NewWithOptions(logoutFunc{logout: func(ctx context.Context, token string) error {
				calls++
				if token != "secret-refresh" {
					t.Fatalf("wrong refresh token: %q", token)
				}
				return tc.err
			}}, Options{CookieSecure: tc.secure, AllowedOrigin: tc.allowedOrigin})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			req.Header.Set("Origin", tc.origin)
			req.Header.Set("Sec-Fetch-Site", tc.fetchSite)
			name := "refresh_token"
			if tc.secure {
				name = "__Host-refresh_token"
			}
			if tc.cookie {
				req.AddCookie(&http.Cookie{Name: name, Value: "secret-refresh"})
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.status || (calls == 1) != tc.called || calls > 1 {
				t.Fatalf("status=%d calls=%d body=%s", w.Code, calls, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Pragma") != "no-cache" {
				t.Fatal("logout response can be cached")
			}
			if strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "secret-refresh") {
				t.Fatal("response leaked secrets")
			}
			if tc.status == 204 && w.Body.Len() != 0 {
				t.Fatal("204 response has a body")
			}
			cookies := w.Result().Cookies()
			if !tc.cleared {
				if len(cookies) != 0 {
					t.Fatal("unexpected cookie deletion")
				}
				return
			}
			if len(cookies) != 1 {
				t.Fatalf("expected one deleted cookie, got %d", len(cookies))
			}
			c := cookies[0]
			if c.Name != name || c.Value != "" || c.Path != "/" || c.Domain != "" || !c.HttpOnly || c.Secure != tc.secure || c.SameSite != http.SameSiteStrictMode || c.MaxAge != -1 || c.Expires.IsZero() || !c.Expires.Before(time.Now()) {
				t.Fatalf("incorrect deletion cookie: %v", c)
			}
		})
	}
}
