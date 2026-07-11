package handler

import "diceDasher/pkg/httputil"

func RegisterRouters(r *httputil.Router) {
	r.Handle("/register").POST()
	r.Handle("/login").POST()
	r.Handle("/logout").POST()
	r.Handle("/me").GET()
	r.Handle("/user").PATCH()
	r.Handle("/password").PUT()
	r.Handle("/user").DELETE()
}
