package handler

import (
	"backend/internal/user"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
	"time"
)

func (h *Handler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var req updateUserRequest

	if err := httputil.UnpackJSON(r, &req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		return
	}

	claims, ok := AccessClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	err := h.user.UpdateUser(r.Context(), user.UserUpdateRecord{ID: claims.UserID, Username: req.Username, Email: req.Email})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrAlreadyExists):
			http.Error(w, "username or email already exists", http.StatusConflict)
		case errors.Is(err, user.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			logger.LogErrorf(r.Context(), "updating user: %s", err)
			http.Error(w, "failed to update user", http.StatusInternalServerError)
		}
		return
	}

	if err := httputil.PackJSON(w, http.StatusOK, "success"); err != nil {
		logger.LogErrorf(r.Context(), "writing user update response: %s", err)
	}

}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := AccessClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	err := h.user.DeleteUser(r.Context(), claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrAlreadyExists):
			http.Error(w, "username or email already exists", http.StatusConflict)
		case errors.Is(err, user.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			logger.LogErrorf(r.Context(), "updating user: %s", err)
			http.Error(w, "failed to update user", http.StatusInternalServerError)
		}
		return
	}

	name := "refresh_token"
	if h.options.CookieSecure {
		name = "__Host-refresh_token"
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
	}) //access revoked
	if err := httputil.PackJSON(w, http.StatusOK, "success"); err != nil {
		logger.LogErrorf(r.Context(), "writing user delete response: %s", err)
	}
}
