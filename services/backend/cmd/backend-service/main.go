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
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		logger.Logf(context.Background(), "backend failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	var shutdownErr error

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	repo, err := dbutil.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer repo.Close()

	tokens, err := auth.NewJWTManager(cfg.JWTKey, cfg.JWTIssuer, cfg.JWTAudience)
	if err != nil {
		return fmt.Errorf("create JWT manager: %w", err)
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
		return fmt.Errorf("listen: %w", err)
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
			err = fmt.Errorf("serve HTTP: %w", err)
			return err
		}
		return nil

	case <-stopCtx.Done():
		logger.Logf(context.Background(), "backend service stopping")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		shutdownErr = fmt.Errorf("graceful shutdown: %w", err)

		if err := srv.Close(); err != nil {
			shutdownErr = errors.Join(shutdownErr, fmt.Errorf("close HTTP server: %w", err))
		}
	}

	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		shutdownErr = errors.Join(shutdownErr, fmt.Errorf("serve HTTP: %w", err))
	}

	if shutdownErr != nil {
		return shutdownErr
	}

	logger.Logf(context.Background(), "backend service shutdown")
	return nil
}
