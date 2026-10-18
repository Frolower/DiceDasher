package handler

import (
	"diceDasher/pkg/httputil"
	"net/http"
)

func (h *Handler) RegisterRouters(r *httputil.Router) {
	// auth specific routes
	r.Handle("/register").POST(h.HandleCreateUser)
	r.Handle("/login").POST(h.HandleLogin)
	r.Handle("/logout").POST(h.RequreAuth(h.HandleLogout))
	r.Handle("/me").GET(h.RequreAuth(notImplemented))
	r.Handle("/user").PATCH(h.RequreAuth(notImplemented)).DELETE(h.RequreAuth(notImplemented))
	r.Handle("/password").PUT(h.RequreAuth(notImplemented))
	r.Handle("/refresh").POST(h.HandleRefresh)

	//health
	r.Handle("/health").GET(Health)
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
