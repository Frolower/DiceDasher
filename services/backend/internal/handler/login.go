package handler

import (
	"backend/internal/auth"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"mime"
	"net/http"
	"time"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	// verifying origin
	if origin := r.Header.Get("Origin"); origin != "" && (h.options.AllowedOrigin == "" || origin != h.options.AllowedOrigin) {
		http.Error(w, "origin is not allowed", http.StatusForbidden)
		return
	}
	if r.Header.Get("Origin") == "" && r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		http.Error(w, "cross-site login is not allowed", http.StatusForbidden)
		return
	}
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" {
		http.Error(w, "application/json is required", http.StatusUnsupportedMediaType)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var req loginRequest
	if err := httputil.UnpackJSON(r, &req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		return
	}
	tokens, err := h.authenticator.Login(r.Context(), auth.LoginInput{Username: req.Username, Password: req.Password, UserAgent: r.UserAgent()})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, auth.ErrUnableToLogin):
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
		default:
			logger.Logf(r.Context(), "ERROR: login failed internally")
			http.Error(w, "failed to log in", http.StatusInternalServerError)
		}
		return
	}
	h.setRefreshCookie(w, tokens)
	if err := httputil.PackJSON(w, http.StatusOK, tokenResponse{AccessToken: tokens.AccessToken, TokenType: "Bearer", ExpiresAt: tokens.AccessExpiresAt}); err != nil {
		logger.Logf(r.Context(), "ERROR writing login response: %s", err)
	}
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	if origin := r.Header.Get("Origin"); origin != "" && (h.options.AllowedOrigin == "" || origin != h.options.AllowedOrigin) {
		http.Error(w, "origin is not allowed", http.StatusForbidden)
		return
	}
	if r.Header.Get("Origin") == "" && r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		http.Error(w, "cross-site logout is not allowed", http.StatusForbidden)
		return
	}

	name := "refresh_token"
	if h.options.CookieSecure {
		name = "__Host-refresh_token"
	}

	cookie, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "invalid cookie", http.StatusBadRequest)
		return
	}

	if cookie != nil {
		if err := h.authenticator.Logout(r.Context(), cookie.Value); err != nil {
			if !errors.Is(err, auth.ErrSessionNotFound) {
				logger.Logf(r.Context(), "ERROR logout failed internally")
				http.Error(w, "failed to logout", http.StatusInternalServerError)
				return
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.options.CookieSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}
