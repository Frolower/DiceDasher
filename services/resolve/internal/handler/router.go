package handler

import "diceDasher/pkg/httputil"

func RegisterRoutes(r *httputil.Router) {
	r.Handle("/resolve").POST(ResolveRoll)

	//health
	r.Handle("/health").GET(Health)
}
