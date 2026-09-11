package main

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository"
	"context"
	"diceDasher/pkg/dbutil"
	"diceDasher/pkg/httputil"
	"log"
	"net"
	"net/http"
	"time"
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

	registration := auth.NewService(repository.New(repo.Pool()))
	r := httputil.NewRouter()
	handler.New(registration).RegisterRouters(r)

	// Wrap with middlewares
	wrapped := httputil.RequestLoggerWithMode(r, cfg.LogMode)
	wrapped = httputil.CORS("http://localhost:8082")(wrapped)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           wrapped,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("backend service READY on %s", cfg.HTTPAddr)
	log.Fatal(srv.Serve(ln))
}
