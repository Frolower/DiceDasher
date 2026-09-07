package handler

import (
	"context"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"diceDasher/services/character/internal/service"
	"diceDasher/services/character/internal/system"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type CharacterCreator interface {
	Create(context.Context, string, json.RawMessage) (uuid.UUID, error)
}

type Handler struct{ characters CharacterCreator }

func New(characters CharacterCreator) *Handler { return &Handler{characters: characters} }

type createCharacterResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handler) postUserCreatedCharacterHandler(w http.ResponseWriter, r *http.Request) {
	sys := r.URL.Query().Get("system")
	if sys == "" {
		http.Error(w, "missing query param: system", http.StatusBadRequest)
		return
	}
	var raw json.RawMessage
	if err := httputil.UnpackJSON(r, &raw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.characters.Create(r.Context(), sys, raw)
	if err != nil {
		logger.Logf(r.Context(), "ERROR creating character: %s", err)
		var input *system.InputError
		var validation *system.ValidationError
		var save *service.SaveError
		switch {
		// Persistence errors must not be mistaken for user input errors in their cause chain.
		case errors.As(err, &save):
			http.Error(w, "failed to save character", http.StatusInternalServerError)
		case errors.Is(err, service.ErrUnknownSystem):
			http.Error(w, "unknown system", http.StatusNotFound)
		case errors.As(err, &input):
			http.Error(w, input.Error(), http.StatusBadRequest)
		case errors.As(err, &validation):
			if err := httputil.PackJSON(w, http.StatusUnprocessableEntity, validation); err != nil {
				logger.Logf(r.Context(), "ERROR encoding validation response: %s", err)
			}
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	if err := httputil.PackJSON(w, http.StatusCreated, createCharacterResponse{ID: id}); err != nil {
		logger.Logf(r.Context(), "ERROR encoding creation response: %s", err)
	}
}
