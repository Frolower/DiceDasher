package handler

import (
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	claims, ok := AccessClaimsFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var request struct {
		System string `json:"system"`
	}
	if err := httputil.UnpackJSON(r, &request); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		return
	}
	if strings.TrimSpace(request.System) == "" {
		http.Error(w, "system is required", http.StatusBadRequest)
		return
	}
	if h.rooms == nil {
		http.Error(w, "room service unavailable", http.StatusServiceUnavailable)
		return
	}
	roomID, err := h.rooms.Create(r.Context(), request.System, claims.UserID)
	if err != nil {
		logger.LogErrorf(r.Context(), "creating room: %v", err)
		switch status.Code(err) {
		case codes.InvalidArgument:
			http.Error(w, "invalid room request", http.StatusBadRequest)
		case codes.Unavailable:
			http.Error(w, "room service unavailable", http.StatusServiceUnavailable)
		case codes.DeadlineExceeded:
			http.Error(w, "room service timed out", http.StatusGatewayTimeout)
		default:
			http.Error(w, "room service failed", http.StatusBadGateway)
		}
		return
	}
	if err := httputil.PackJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
		RoomID string `json:"room_id"`
	}{Status: "ok", RoomID: roomID}); err != nil {
		logger.LogErrorf(r.Context(), "writing room response: %v", err)
	}
}
