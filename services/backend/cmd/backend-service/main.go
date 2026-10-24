package main

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/user"
	"context"
	"diceDasher/pkg/dbutil"
	"diceDasher/pkg/httputil"
	"diceDasher/pkg/logger"
	"net"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.LogFatalf(context.Background(), "%v", err)
	}

	repo, err := dbutil.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.LogFatalf(context.Background(), "%v", err)
	}
	defer repo.Close()

	tokens, err := auth.NewJWTManager(cfg.JWTKey, cfg.JWTIssuer, cfg.JWTAudience)
	if err != nil {
		logger.LogFatalf(context.Background(), "%v", err)
	}
	store := repository.New(repo.Pool())
	// A single repository implements both the user and session interfaces.
	authentication := auth.NewServiceWithDependencies(store, store, auth.NewBcryptHasher(), tokens)
	users := user.NewServiceWithDependencies(store, user.NewBcryptHasher())
	r := httputil.NewRouter()
	handler.NewWithOptions(authentication, users, handler.Options{CookieSecure: cfg.CookieSecure, AllowedOrigin: cfg.AllowedOrigin}).RegisterRouters(r)

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
		logger.LogFatalf(context.Background(), "%v", err)
	}

	logger.Logf(context.Background(), "backend service READY on %s", cfg.HTTPAddr)
	logger.LogFatalf(context.Background(), "%v", srv.Serve(ln))
}
