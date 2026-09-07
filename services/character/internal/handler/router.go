package handler

import "diceDasher/pkg/httputil"

func (h *Handler) RegisterRouters(r *httputil.Router) {
	r.Handle("/character").POST(h.postUserCreatedCharacterHandler)

	//health
	r.Handle("/health").GET(Health)
}
