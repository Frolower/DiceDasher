package handler

import (
	"log/slog"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		slog.Warn("health write failed", "err", err)
	}
}
