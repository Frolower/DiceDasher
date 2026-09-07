package handler

import "diceDasher/pkg/httputil"

func (h *Handler) RegisterRoutes(r *httputil.Router) {
	r.Handle("/resolve").POST(h.Resolve) // POST /resolve?system=system-name&action=action-name

	//health
	r.Handle("/health").GET(Health)
}
