package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"diceDasher/pkg/logger"
	"diceDasher/services/room/internal/config"
	"diceDasher/services/room/internal/server"
)

func main() {
	if err := run(); err != nil {
		logger.LogFatalf(context.Background(), "%v", err)
	}
}

func run() error {
	cfg := config.Load()
	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}
	defer listener.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New()
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(listener) }()
	logger.Logf(ctx, "room service READY on %s", listener.Addr())

	select {
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("serve gRPC: %w", err)
		}
		return nil
	case <-ctx.Done():
		logger.Logf(context.Background(), "room service stopping")
	}
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown gRPC: %w", err)
	}
	if err := <-serveErr; err != nil {
		return fmt.Errorf("serve gRPC: %w", err)
	}
	logger.Logf(context.Background(), "room service shutdown")
	return nil
}
