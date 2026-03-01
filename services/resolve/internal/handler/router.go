package handler

import "diceDasher/pkg/httputil"

func RegisterRoutes(r *httputil.Router) {
	r.Handle("/resolve/basic").POST(ResolveRoll)
	r.Handle("/resolve").POST(ResolveSystem) // POST /resolve?system=system-name

	//health
	r.Handle("/health").GET(Health)
}
