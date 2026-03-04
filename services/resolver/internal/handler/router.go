package handler

import "diceDasher/pkg/httputil"

func RegisterRoutes(r *httputil.Router) {
	r.Handle("/resolver").POST(ResolveHandler) // POST /resolver?system=system-name&action=action-name

	//health
	r.Handle("/health").GET(Health)
}
