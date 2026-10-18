package handler

import (
	"backend/internal/auth"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"errors"
	"net/http"
)

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Limit the body before decoding. 16 KiB is enough for three fields;
	// an oversized request must not consume unbounded memory.
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	var req createUserRequest
	// Decode directly into a struct: UnpackJSON rejects unknown fields,
	// invalid JSON, and multiple JSON values in one body.
	if err := httputil.UnpackJSON(r, &req); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			// Do not return the decoder error text: it may contain submitted values.
			http.Error(w, "invalid request body", http.StatusBadRequest)
		}
		return
	}

	// Convert the HTTP model into input for the registration flow.
	// Field validation, hashing, and persistence happen inside auth.Service.
	_, err := h.authenticator.Register(r.Context(), auth.CreateInput{
		Username: req.Username, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, auth.ErrAlreadyExists):
			http.Error(w, "username or email already exists", http.StatusConflict)
		default:
			// Do not expose database details to the client. The repository also excludes
			// PostgreSQL Detail from errors because it may contain the user's email.
			logger.LogErrorf(r.Context(), "registering user: %s", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
		}
		return
	}

	// 201 means the user has been saved. Registration does not yet create
	// a session, cookies, or tokens; the internal UUID does not need to be exposed.
	if err := httputil.PackJSON(w, http.StatusCreated, createUserResponse{Status: "created"}); err != nil {
		// Headers have already been sent: do not write a second HTTP response here.
		logger.LogErrorf(r.Context(), "writing registration response: %s", err)
	}
}
