package server

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	reflectionpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/test/bufconn"
)

func startServer(t *testing.T) (*Server, *grpc.ClientConn, <-chan error) {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	srv := New()
	done := make(chan error, 1)
	go func() { done <- srv.Serve(listener) }()
	t.Cleanup(func() {
		srv.grpc.Stop()
		listener.Close()
	})
	conn, err := grpc.NewClient("passthrough:///room",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return srv, conn, done
}

func TestHealthReflectionAndShutdown(t *testing.T) {
	srv, conn, done := startServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if response.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("health = %s, want SERVING", response.GetStatus())
	}

	reflectionCtx, stopReflection := context.WithCancel(ctx)
	defer stopReflection()
	stream, err := reflectionpb.NewServerReflectionClient(conn).ServerReflectionInfo(reflectionCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := stream.Send(&reflectionpb.ServerReflectionRequest{
		MessageRequest: &reflectionpb.ServerReflectionRequest_ListServices{ListServices: ""},
	}); err != nil {
		t.Fatal(err)
	}
	services, err := stream.Recv()
	if err != nil {
		t.Fatal(err)
	}
	foundHealth := false
	for _, service := range services.GetListServicesResponse().GetService() {
		if service.GetName() == "grpc.health.v1.Health" {
			foundHealth = true
		}
	}
	if !foundHealth {
		t.Fatal("health service missing from reflection")
	}
	stopReflection()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Serve after shutdown: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Serve did not return after shutdown")
	}
	response, err = srv.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil || response.GetStatus() != healthpb.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("health after shutdown = %v, %v; want NOT_SERVING", response, err)
	}
}

func TestShutdownDeadlineStopsActiveWatch(t *testing.T) {
	srv, conn, done := startServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	watch, err := healthpb.NewHealthClient(conn).Watch(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := watch.Recv(); err != nil {
		t.Fatal(err)
	}

	shutdownCtx, stop := context.WithTimeout(ctx, 50*time.Millisecond)
	defer stop()
	if err := srv.Shutdown(shutdownCtx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown = %v, want deadline exceeded for active stream", err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("Serve did not return after forced shutdown")
	}
}
