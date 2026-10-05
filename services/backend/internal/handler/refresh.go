package handler

import (
	"backend/internal/auth"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
)

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	if origin := r.Header.Get("Origin"); origin != "" && (h.options.AllowedOrigin == "" || origin != h.options.AllowedOrigin) {
		http.Error(w, "origin is not allowed", http.StatusForbidden)
		return
	}
	if r.Header.Get("Origin") == "" && r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		http.Error(w, "cross-site refresh is not allowed", http.StatusForbidden)
		return
	}

	name := "refresh_token"
	if h.options.CookieSecure {
		name = "__Host-refresh_token"
	}

	cookie, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "invalid cookie", http.StatusBadRequest)
		return
	}

	tokens, err := h.authenticator.Refresh(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, auth.ErrSessionNotFound) {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		} else {
			logger.Logf(r.Context(), "ERROR refresh failed internally")
			http.Error(w, "failed to refresh the token", http.StatusInternalServerError)
		}
		return
	}

	h.setRefreshCookie(w, tokens)

	if err := httputil.PackJSON(w, http.StatusOK, tokenResponse{
		AccessToken: tokens.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   tokens.AccessExpiresAt,
	}); err != nil {
		logger.Logf(r.Context(), "ERROR writing refresh response: %s", err)
	}
}
