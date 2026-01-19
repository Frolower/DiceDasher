package main

import (
	"diceDasher/services/resolve/internal/config"
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

	r := httputil.NewRouter()
	handler.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}
