package handler

import (
	"backend/internal/auth"
	"context"
	"diceDasher/pkg/httputil"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type loginFunc struct {
	Authenticator
	login func(context.Context, auth.LoginInput) (auth.Tokens, error)
}

func (f loginFunc) Login(c context.Context, i auth.LoginInput) (auth.Tokens, error) {
	return f.login(c, i)
}

func TestLoginHTTP(t *testing.T) {
	for _, tc := range []struct {
		name, body, origin, contentType string
		err                             error
		status                          int
		called                          bool
	}{
		{name: "success", body: `{"username":"Alice","password":"StrongPass1"}`, status: 200, called: true},
		{name: "wrong credentials", body: `{}`, err: auth.ErrUnableToLogin, status: 401, called: true},
		{name: "invalid input", body: `{}`, err: auth.InvalidInputError{Err: errors.New("missing fields")}, status: 400, called: true},
		{name: "database failure", body: `{}`, err: errors.New("private password/hash details"), status: 500, called: true},
		{name: "malformed", body: `{`, status: 400},
		{name: "unknown field", body: `{"user_agent":"fake"}`, status: 400},
		{name: "wrong type", body: `{"password":12}`, status: 400},
		{name: "large", body: `{"password":"` + strings.Repeat("x", 17000) + `"}`, status: 413},
		{name: "foreign origin", body: `{}`, origin: "https://evil.example", status: 403},
		{name: "opaque origin", body: `{}`, origin: "null", status: 403},
		{name: "form CSRF", body: `username=Alice`, contentType: "application/x-www-form-urlencoded", status: 415},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			now := time.Now().UTC().Truncate(time.Second)
			h := NewWithOptions(loginFunc{login: func(_ context.Context, input auth.LoginInput) (auth.Tokens, error) {
				called = true
				if input.UserAgent != "real agent" {
					t.Fatal("lost User-Agent header")
				}
				if tc.name == "success" && (input.Username != "Alice" || input.Password != "StrongPass1") {
					t.Fatal("lost credentials")
				}
				return auth.Tokens{AccessToken: "access", RefreshToken: "secret-refresh", AccessExpiresAt: now.Add(time.Minute), RefreshExpiresAt: now.Add(time.Hour)}, tc.err
			}}, Options{CookieSecure: true, AllowedOrigin: "https://app.example"})
			router := httputil.NewRouter()
			h.RegisterRouters(router)
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tc.body))
			ct := tc.contentType
			if ct == "" {
				ct = "application/json"
			}
			req.Header.Set("Content-Type", ct)
			req.Header.Set("User-Agent", "real agent")
			origin := tc.origin
			if origin == "" {
				origin = "https://app.example"
			}
			req.Header.Set("Origin", origin)
			w := httptest.NewRecorder()
			CORS("https://app.example")(router).ServeHTTP(w, req)
			if w.Code != tc.status || called != tc.called {
				t.Fatalf("code=%d called=%v body=%s", w.Code, called, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" || strings.Contains(w.Body.String(), "private") || strings.Contains(w.Body.String(), "secret-refresh") {
				t.Fatal("response leaked secret or is cacheable")
			}
			cookies := w.Result().Cookies()
			if tc.status != 200 {
				if len(cookies) != 0 {
					t.Fatal("cookie on failed login")
				}
				return
			}
			if len(cookies) != 1 {
				t.Fatal("missing refresh cookie")
			}
			c := cookies[0]
			if c.Name != "__Host-refresh_token" || c.Value != "secret-refresh" || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteStrictMode || c.Path != "/" || c.Domain != "" || c.MaxAge <= 0 || !c.Expires.Equal(now.Add(time.Hour)) {
				t.Fatalf("incorrect cookie: %v", c)
			}
			var response tokenResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.AccessToken != "access" || response.TokenType != "Bearer" {
				t.Fatalf("invalid response: %s", w.Body.String())
			}
			if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
				t.Fatal("missing credential CORS")
			}
		})
	}
}

func TestLoginLocalCookie(t *testing.T) {
	h := NewWithOptions(nil, Options{CookieSecure: false})
	w := httptest.NewRecorder()
	h.setRefreshCookie(w, auth.Tokens{RefreshToken: "refresh", RefreshExpiresAt: time.Now().Add(time.Hour)})
	c := w.Result().Cookies()[0]
	if c.Secure || !c.HttpOnly || c.Name != "refresh_token" {
		t.Fatal("incorrect local cookie")
	}
}
