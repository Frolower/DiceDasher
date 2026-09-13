package handler

import (
	"backend/internal/auth"
	"net/http"
	"time"
)

func (h *Handler) setRefreshCookie(w http.ResponseWriter, tokens auth.Tokens) {
	name := "refresh_token"
	// __Host- не позволяет расширить cookie на соседние поддомены.
	// Локальный HTTP использует другое имя, потому что __Host- требует Secure.
	if h.options.CookieSecure {
		name = "__Host-refresh_token"
	}
	maxAge := int(time.Until(tokens.RefreshExpiresAt).Seconds())
	if maxAge < 1 {
		maxAge = 1
	}
	http.SetCookie(w, &http.Cookie{
		Name: name, Value: tokens.RefreshToken, Path: "/", HttpOnly: true,
		Secure: h.options.CookieSecure, SameSite: http.SameSiteStrictMode,
		Expires: tokens.RefreshExpiresAt, MaxAge: maxAge,
	})
}
