package handler

import (
	"backend/internal/auth"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
	"strings"
)

func (h *Handler) RequreAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		headers := r.Header.Values("Authorization")
		if len(headers) != 1 {
			writeUnauthorized(w)
			return
		}
		parts := strings.Fields(headers[0])
		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			writeUnauthorized(w)
			return
		}

		claims, err := h.authenticator.Authenticate(
			r.Context(),
			parts[1],
		)
		if err != nil {
			switch {
			case errors.Is(err, auth.ErrInvalidToken),
				errors.Is(err, auth.ErrSessionNotFound):
				writeUnauthorized(w)
			default:
				logger.Logf(
					r.Context(),
					"ERROR authentication failed internally",
				)
				http.Error(
					w,
					"authentication failed",
					http.StatusInternalServerError,
				)
			}
			return
		}

		ctx := withAccessClaims(r.Context(), claims)
		next(w, r.WithContext(ctx))
	}
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="backend"`)
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}