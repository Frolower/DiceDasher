package handler

import (
	"backend/internal/user"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
)

func (h *Handler) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var req updatePasswordRequest

	claims, ok := AccessClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	if err := httputil.UnpackJSON(r, &req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		return
	}

	err := h.user.UpdatePassword(r.Context(), user.PasswordUpdateRecord{Id: claims.UserID, OldPassword: req.OldPassword, NewPassword: req.NewPassword})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		case errors.Is(err, user.ErrOldPasswordDoesntMatch):
			http.Error(w, "old password does not match", http.StatusBadRequest)
		default:
			logger.LogErrorf(r.Context(), "updating user: %s", err)
			http.Error(w, "failed to update password", http.StatusInternalServerError)
		}
		return
	}

	if err := httputil.PackJSON(w, http.StatusOK, "success"); err != nil {
		logger.LogErrorf(r.Context(), "writing password update response: %s", err)
	}
}
