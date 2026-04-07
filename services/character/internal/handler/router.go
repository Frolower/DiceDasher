package handler

import "diceDasher/pkg/httputil"

func RegisterRouters(r *httputil.Router) {
	r.Handle("/character").POST(postUserCreatedCharacterHandler)

	//health
	r.Handle("/health").GET(Health)
}
