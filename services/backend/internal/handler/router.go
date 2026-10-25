package handler

import (
	"diceDasher/pkg/httputil"
	"net/http"
)

func (h *Handler) RegisterRouters(r *httputil.Router) {
	// auth specific routes
	r.Handle("/register").POST(h.HandleCreateUser)
	r.Handle("/login").POST(h.HandleLogin)
	r.Handle("/logout").POST(h.HandleLogout)
	r.Handle("/me").GET(h.RequreAuth(h.HandleMe))
	r.Handle("/refresh").POST(h.HandleRefresh)

	// user specific routes
	r.Handle("/user").PATCH(h.RequreAuth(h.HandleUpdateUser)).DELETE(h.RequreAuth(h.HandleDeleteUser))
	r.Handle("/password").PUT(h.RequreAuth(notImplemented))

	//health
	r.Handle("/health").GET(Health)
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
