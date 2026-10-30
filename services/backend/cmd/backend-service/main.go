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
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	handler.NewWithOptions(authentication, users, handler.Options{CookieSecure: cfg.CookieSecure, AllowedOrigins: cfg.AllowedOrigins}).RegisterRouters(r)

	// Wrap with middlewares
	wrapped := httputil.RequestLoggerWithMode(r, cfg.LogMode)
	wrapped = handler.CORS(cfg.AllowedOrigins)(wrapped)

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

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)

	go func() {
		serveErr <- srv.Serve(ln)
	}()

	logger.Logf(context.Background(), "backend service READY on %s", cfg.HTTPAddr)

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logf(context.Background(), "HTTP server failed: %v", err)
		}
		return

	case <-stopCtx.Done():
		logger.Logf(context.Background(), "backend service stopping")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logf(context.Background(), "graceful shutdown failed: %v", err)

		if err := srv.Close(); err != nil {
			logger.Logf(context.Background(), "close HTTP server: %v", err)
		}
	}

	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Logf(context.Background(), "HTTP server failed: %v", err)
	}

	logger.Logf(context.Background(), "backend service shutdown")
}
