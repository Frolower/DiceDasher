package handler

import (
	"diceDasher/pkg/logger"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		logger.LogWarningf(r.Context(), "health write failed: %v", err)
	}
}
