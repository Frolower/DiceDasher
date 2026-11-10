package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpc   *grpc.Server
	health *health.Server
}

func New() *Server {
	s := &Server{grpc: grpc.NewServer(), health: health.NewServer()}
	healthpb.RegisterHealthServer(s.grpc, s.health)
	reflection.Register(s.grpc)
	// Register future room RPC services here.
	return s
}

func (s *Server) Serve(listener net.Listener) error {
	return s.grpc.Serve(listener)
}

// Shutdown reports NOT_SERVING, drains RPCs, and forces a stop if ctx expires.
func (s *Server) Shutdown(ctx context.Context) error {
	s.health.Shutdown()
	done := make(chan struct{})
	go func() {
		s.grpc.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.grpc.Stop()
		<-done
		return ctx.Err()
	}
}
