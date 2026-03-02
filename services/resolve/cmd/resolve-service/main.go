package main

import (
	"diceDasher/services/resolve/internal/config"
	"diceDasher/services/resolve/internal/system"
	"diceDasher/services/resolve/internal/system/generic"
	"diceDasher/services/resolve/internal/system/tes"
	"diceDasher/services/resolve/internal/system/vtmv5"
	"log"
	"net/http"

	"diceDasher/pkg/httputil"
	"diceDasher/services/resolve/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// register resolvers
	system.Register("generic", generic.Resolver{})
	system.Register("tes", tes.Resolver{})
	system.Register("vtmv5", vtmv5.Resolver{})

	r := httputil.NewRouter()
	handler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}
