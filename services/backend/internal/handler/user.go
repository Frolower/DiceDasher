package handler

import (
	"backend/internal/repository"
	"backend/internal/user"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	var raw json.RawMessage

	if err := httputil.UnpackJSON(r, &raw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Logf(r.Context(), "error unpacking request body: %s", err)
		return
	}

	if err := json.Unmarshal(raw, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Logf(r.Context(), "error decoding create user request: %s", err)
		return
	}

	repo, err := repository.FromContext(r.Context())
	if err != nil {
		logger.Logf(r.Context(), "ERROR: repository not in context %s", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	service := user.NewService(repo)
	created, err := service.Create(r.Context(), user.CreateInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, user.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Logf(r.Context(), "ERROR validating user: %s", err)
			return
		}

		logger.Logf(r.Context(), "ERROR creating user: %s", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	if err := httputil.PackJSON(w, http.StatusCreated, createUserResponse{ID: created.ID}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
