package handler

import "diceDasher/pkg/httputil"

func RegisterRouters(r *httputil.Router) {
	// user specific routes
	r.Handle("/register").POST()
	r.Handle("/login").POST()
	r.Handle("/logout").POST()
	r.Handle("/me").GET()
	r.Handle("/user").PATCH().DELETE()
	r.Handle("/password").PUT()

	//health
	r.Handle("/health").GET(health)
}
