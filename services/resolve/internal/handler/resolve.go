package handler

import (
	"diceDasher/services/resolve/internal/validation"
	"fmt"
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

func ResolveSystem(w http.ResponseWriter, r *http.Request) {
	sys := r.URL.Query().Get("system")
	if sys == "" {
		http.Error(w, "missing query param: system", http.StatusBadRequest)
		return
	}

	switch sys {
	case "tes":
		resolve := model.ResolveTES{}

		if err := httputil.UnpackJSON(r, &resolve); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validation.ValidateResolveTES(resolve); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		response := service.ResolveTES(resolve)

		if err := httputil.PackJSON(w, http.StatusOK, response); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		return

	case "vtmV5":
		fmt.Println("WIP")
		return
	default:
		http.Error(w, "this system is not implemented", http.StatusNotFound)
		return
	}
}
