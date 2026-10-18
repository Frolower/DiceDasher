package handler

import (
	"backend/internal/auth"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
)

func (h *Handler) HandleMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := AccessClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	profile, err := h.authenticator.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			writeUnauthorized(w)
		default:
			logger.LogErrorf(r.Context(), "loading user profile")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	if err := httputil.PackJSON(w, http.StatusOK, meResponse{
		ID:       profile.ID,
		Username: profile.Username,
		Email:    profile.Email}); err != nil {
		logger.LogErrorf(r.Context(), "writing profile response: %s", err)
	}
}
