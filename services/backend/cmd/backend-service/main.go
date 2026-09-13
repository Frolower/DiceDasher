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

	tokens, err := auth.NewJWTManager(cfg.JWTKey, cfg.JWTIssuer, cfg.JWTAudience)
	if err != nil {
		log.Fatal(err)
	}
	store := repository.New(repo.Pool())
	// Один repository реализует интерфейсы пользователей и сессий.
	authentication := auth.NewServiceWithDependencies(store, store, auth.NewBcryptHasher(), tokens)
	r := httputil.NewRouter()
	handler.NewWithOptions(authentication, handler.Options{CookieSecure: cfg.CookieSecure, AllowedOrigin: cfg.AllowedOrigin}).RegisterRouters(r)

	// Wrap with middlewares
	wrapped := httputil.RequestLoggerWithMode(r, cfg.LogMode)
	wrapped = handler.CORS(cfg.AllowedOrigin)(wrapped)

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
