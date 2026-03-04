package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/system"
)

func ResolveHandler(repo *repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sys := r.URL.Query().Get("system")
		if sys == "" {
			http.Error(w, "missing query param: system", http.StatusBadRequest)
			return
		}

		action := r.URL.Query().Get("action")
		if action == "" {
			action = "roll" // default action
			log.Println("request arrived without action")
		}

		resolver, err := system.Get(sys)
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
			return
		}

		resp, status, err := resolver.Resolve(r.Context(), action, raw)
		if err != nil {
			http.Error(w, err.Error(), status)
			logger.Logf(r.Context(), "ERROR: %s", err)
			return
		}

		// Save roll history to database
		_, err = repo.InsertRollHistory(r.Context(), repository.RollHistory{
			SystemName:      sys,
			ActionType:      action,
			RequestPayload:  raw,
			ResponsePayload: resp,
		})
		if err != nil {
			logger.Logf(r.Context(), "ERROR saving roll history: %s", err)
		}

		if err := httputil.PackJSON(w, status, resp); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}
