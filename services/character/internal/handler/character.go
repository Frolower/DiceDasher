package handler

import (
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

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

	resp, status, err := character.CreateCharacter(r.Context(), raw)
	if err != nil {
		http.Error(w, err.Error(), status)
		logger.Logf(r.Context(), "ERROR: %s", err)
		return
	}
}
