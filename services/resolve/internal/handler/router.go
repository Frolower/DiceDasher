package handler

import "diceDasher/pkg/httputil"

func RegisterRoutes(r *httputil.Router) {
	r.Handle("/resolve").POST(Resolve) // POST /resolve?system=system-name

	//health
	r.Handle("/health").GET(Health)
}
