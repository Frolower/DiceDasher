package main

import (
	"context"
	"diceDaher/service/character/internal/config"
	"diceDaher/service/character/internal/handler"
	"diceDasher/pkg/dbutil"
	"diceDasher/pkg/httputil"
	"log"
	"net"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := dbutil.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	r := httputil.NewRouter()
	handler.RegisterRouters(r)

	// Wrap with middlewares
	wrapped := handler.WithRepository(repo)(r)
	wrapped = httputil.RequestLoggerWithMode(wrapped, "default")

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: wrapped,
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("character service READY on %s", cfg.HTTPAddr)
	log.Fatal(srv.Serve(ln))
}
