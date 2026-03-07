package handler

import (
	"diceDasher/pkg/dbutil"
	"net/http"
)

func WithRepository(repo *dbutil.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := dbutil.WithRepo(r.Context(), repo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
