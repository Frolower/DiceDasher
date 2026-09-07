package handler

import (
	"context"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/service"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"errors"
	"net/http"
)

type ResolveService interface {
	Resolve(context.Context, service.Command) (service.Result, error)
}
type Handler struct{ service ResolveService }

func New(service ResolveService) *Handler { return &Handler{service: service} }

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	var raw json.RawMessage
	if err := httputil.UnpackJSON(r, &raw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	requestID, _ := logger.ReqIDFromContext(r.Context())
	result, err := h.service.Resolve(r.Context(), service.Command{
		System: r.URL.Query().Get("system"), Action: r.URL.Query().Get("action"), Payload: raw, RequestID: requestID,
	})
	if err != nil {
		logger.Logf(r.Context(), "ERROR: %s", err)
		status, message := errorResponse(err)
		http.Error(w, message, status)
		return
	}
	if err := httputil.PackJSON(w, http.StatusOK, result); err != nil {
		logger.Logf(r.Context(), "ERROR writing response: %s", err)
	}
}

func errorResponse(err error) (int, string) {
	var requestError *system.RequestError
	switch {
	case errors.As(err, &requestError):
		switch requestError.Kind {
		case system.BadRequest:
			return http.StatusBadRequest, requestError.Error()
		case system.Validation:
			return http.StatusUnprocessableEntity, requestError.Error()
		}
	case errors.Is(err, system.ErrUnknownSystem):
		return http.StatusNotFound, "unknown system"
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "record not found"
	case errors.Is(err, system.ErrInvalidTransition):
		return http.StatusConflict, system.ErrInvalidTransition.Error()
	case errors.Is(err, system.ErrLegacyContinuation):
		return http.StatusConflict, system.ErrLegacyContinuation.Error()
	}
	return http.StatusInternalServerError, "internal error"
}
