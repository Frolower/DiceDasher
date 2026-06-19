package handler

import (
	"diceDaher/service/character/internal/repository"
	"diceDaher/service/character/internal/system"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type createCharacterResponse struct {
	ID uuid.UUID `json:"id"`
}

func postUserCreatedCharacterHandler(w http.ResponseWriter, r *http.Request) {
	sys := r.URL.Query().Get("system")
	if sys == "" {
		http.Error(w, "missing query param: system", http.StatusBadRequest)
		log.Panicln("request arrived without system")
		return
	}

	character, err := system.Get(sys)
	if err != nil {
		if errors.Is(err, system.ErrUnknownSystem) {
			http.Error(w, "unknown system", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var raw json.RawMessage
	if err := httputil.UnpackJSON(r, &raw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Logf(r.Context(), "ERROR: %s", err)
		return
	}

	created, status, err := character.CreateCharacter(r.Context(), raw)
	if err != nil {
		http.Error(w, err.Error(), status)
		logger.Logf(r.Context(), "ERROR: %s", err)
		return
	}

	repo, err := repository.FromContext(r.Context())
	if err != nil {
		logger.Logf(r.Context(), "ERROR: repository not in context: %s", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	id, err := repo.InsertPlayerCreatedCharacter(r.Context(), repository.Character{
		UserID:        created.UserID,
		SystemName:    sys,
		CharacterType: created.CharacterType,
		Name:          created.Name,
		Data:          created.Data,
	})
	if err != nil {
		logger.Logf(r.Context(), "ERROR saving character: %s", err)
		http.Error(w, "failed to save character", http.StatusInternalServerError)
		return
	}

	if err := httputil.PackJSON(w, http.StatusCreated, createCharacterResponse{ID: id}); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
