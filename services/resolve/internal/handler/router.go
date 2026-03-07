package handler

import "diceDasher/pkg/httputil"

func RegisterRoutes(r *httputil.Router) {
	r.Handle("/resolve").POST(ResolveHandler) // POST /resolve?system=system-name&action=action-name

	//health
	r.Handle("/health").GET(Health)
}
