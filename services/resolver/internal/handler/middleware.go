package handler

import (
	"net/http"

	"diceDasher/services/resolve/internal/repository"
)

func WithRepository(repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := repository.WithRepo(r.Context(), repo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
