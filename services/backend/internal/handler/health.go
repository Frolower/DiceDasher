package handler

import (
	"backend/internal/healthcheck"
	"diceDasher/pkg/logger"
	"encoding/json/v2"
	"errors"
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	service := r.URL.Query().Get("service")
	if service == "" || service == "backend" {
		if _, err := w.Write([]byte("OK")); err != nil {
			logger.LogWarningf(r.Context(), "health write failed: %v", err)
		}
		return
	}
	code, status := http.StatusOK, "SERVING"
	err := healthcheck.ErrUnknownService
	if h.healthChecker != nil {
		err = h.healthChecker.Check(r.Context(), service)
	}
	if errors.Is(err, healthcheck.ErrUnknownService) {
		code, status = http.StatusNotFound, "UNKNOWN_SERVICE"
	} else if err != nil {
		code, status = http.StatusServiceUnavailable, "NOT_SERVING"
		logger.LogWarningf(r.Context(), "health check for %s failed: %v", service, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.MarshalWrite(w, struct {
		Service string `json:"service"`
		Status  string `json:"status"`
	}{service, status}); err != nil {
		logger.LogWarningf(r.Context(), "health write failed: %v", err)
	}
}
