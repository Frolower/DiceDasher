package handler

import "diceDasher/pkg/httputil"

func RegisterRouters(r *httputil.Router) {
	r.Handle("/health").GET(Health)
}
