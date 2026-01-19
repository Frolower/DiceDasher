package handler

import (
	"diceDasher/services/resolve/internal/validation"
	"net/http"

	"diceDasher/pkg/httputil"
	"diceDasher/services/resolve/internal/model"
	"diceDasher/services/resolve/internal/service"
)

func ResolveRoll(w http.ResponseWriter, r *http.Request) {
	var resolve model.Resolve

	if err := httputil.UnpackJSON(r, &resolve); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateResolve(resolve); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response := service.ResolveRoll(resolve)

	if err := httputil.PackJSON(w, http.StatusOK, response); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
