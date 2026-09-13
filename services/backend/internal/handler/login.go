package handler

import (
	"backend/internal/auth"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"mime"
	"net/http"
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
