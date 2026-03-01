package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"diceDasher/pkg/httputil"
	"diceDasher/services/resolve/internal/system"
)

func Resolve(w http.ResponseWriter, r *http.Request) {
	sys := r.URL.Query().Get("system")
	if sys == "" {
		http.Error(w, "missing query param: system", http.StatusBadRequest)
		return
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

	resp, status, err := resolver.Resolve(r.Context(), raw)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	if err := httputil.PackJSON(w, status, resp); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
