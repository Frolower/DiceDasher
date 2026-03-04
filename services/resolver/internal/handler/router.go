package handler

import (
	"diceDasher/pkg/httputil"
	"diceDasher/services/resolve/internal/repository"
)

func RegisterRoutes(r *httputil.Router, repo *repository.Repository) {
	r.Handle("/resolver").POST(ResolveHandler(repo)) // POST /resolver?system=system-name&action=action-name

	//health
	r.Handle("/health").GET(Health)
}
