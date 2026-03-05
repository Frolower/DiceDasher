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

	"github.com/google/uuid"
)

type resolveEnvelope struct {
	RecordID *uuid.UUID `json:"record_id,omitempty"`
	Data     any        `json:"data"`
}

func ResolveHandler(w http.ResponseWriter, r *http.Request) {
	sys := r.URL.Query().Get("system")
	if sys == "" {
		http.Error(w, "missing query param: system", http.StatusBadRequest)
		log.Println("request arrived without system")
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
		logger.Logf(r.Context(), "ERROR: %s", err)
		return
	}

	resp, status, err := resolver.Resolve(r.Context(), action, raw)
	if err != nil {
		http.Error(w, err.Error(), status)
		logger.Logf(r.Context(), "ERROR: %s", err)
		return
	}

	// Save roll history to database
	var recordID *uuid.UUID

	repo, err := repository.FromContext(r.Context())
	if err != nil {
		logger.Logf(r.Context(), "ERROR: repository not in context: %s", err)
	} else {
		respJSON, err := json.Marshal(resp)
		if err != nil {
			logger.Logf(r.Context(), "ERROR marshaling response: %s", err)
		} else {
			requestID, exists := logger.ReqIDFromContext(r.Context())
			if exists != true {
				logger.Logf(r.Context(), "ERROR getting request ID")
			} else {
				id, err := repo.InsertRollHistory(r.Context(), repository.RollHistory{
					RequestID:       requestID,
					SystemName:      sys,
					ActionType:      action,
					RequestPayload:  raw,
					ResponsePayload: respJSON,
				})
				if err != nil {
					logger.Logf(r.Context(), "ERROR saving roll history: %s", err)
				} else {
					recordID = &id
				}
			}
		}
	}

	out := resolveEnvelope{Data: resp, RecordID: recordID}

	if err := httputil.PackJSON(w, status, out); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
