package handler

import (
	"backend/internal/auth"
	"net/http"
	"time"
)

func (h *Handler) setRefreshCookie(w http.ResponseWriter, tokens auth.Tokens) {
	name := "refresh_token"
	// __Host- prevents the cookie from being scoped to sibling subdomains.
	// Local HTTP uses a different name because __Host- requires Secure.
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
