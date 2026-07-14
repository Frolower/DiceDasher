package handler

import (
	"diceDasher/pkg/httputil"
	"net/http"
)

func RegisterRouters(r *httputil.Router) {
	// user specific routes
	r.Handle("/register").POST(HandleCreateUser)
	r.Handle("/login").POST(notImplemented)
	r.Handle("/logout").POST(notImplemented)
	r.Handle("/me").GET(notImplemented)
	r.Handle("/user").PATCH(notImplemented).DELETE(notImplemented)
	r.Handle("/password").PUT(notImplemented)

	//health
	r.Handle("/health").GET(Health)
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
