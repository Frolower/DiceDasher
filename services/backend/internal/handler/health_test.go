package handler

import (
	"backend/internal/healthcheck"
	"diceDasher/pkg/httputil"
	"encoding/json/v2"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestHealthHTTPToGRPC(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	defer server.Stop()
	grpcHealth := health.NewServer()
	healthpb.RegisterHealthServer(server, grpcHealth)
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	t.Cleanup(func() { server.Stop(); <-done })
	checker, err := healthcheck.New(map[string]string{"room": "grpc://" + listener.Addr().String()})
	if err != nil {
		t.Fatal(err)
	}
	defer checker.Close()
	router := httputil.NewRouter()
	NewWithOptions(nil, nil, nil, checker, Options{}).RegisterRouters(router)

	for _, tc := range []struct {
		name, query, status string
		grpcStatus          healthpb.HealthCheckResponse_ServingStatus
		code                int
	}{
		{"self", "", "", healthpb.HealthCheckResponse_SERVING, 200},
		{"explicit self", "?service=backend", "", healthpb.HealthCheckResponse_SERVING, 200},
		{"healthy", "?service=room", "SERVING", healthpb.HealthCheckResponse_SERVING, 200},
		{"not serving", "?service=room", "NOT_SERVING", healthpb.HealthCheckResponse_NOT_SERVING, 503},
		{"unknown health", "?service=room", "NOT_SERVING", healthpb.HealthCheckResponse_UNKNOWN, 503},
		{"recovered", "?service=room", "SERVING", healthpb.HealthCheckResponse_SERVING, 200},
		{"unknown service", "?service=missing", "UNKNOWN_SERVICE", healthpb.HealthCheckResponse_SERVING, 404},
		{"URL is not a service", "?service=http://example.com", "UNKNOWN_SERVICE", healthpb.HealthCheckResponse_SERVING, 404},
		{"unreachable", "?service=room", "NOT_SERVING", healthpb.HealthCheckResponse_SERVING, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			grpcHealth.SetServingStatus("", tc.grpcStatus)
			if tc.name == "unreachable" {
				server.Stop()
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health"+tc.query, nil))
			if response.Code != tc.code {
				t.Fatalf("HTTP = %d, want %d: %s", response.Code, tc.code, response.Body.String())
			}
			if response.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("health response must not be cached")
			}
			if tc.status == "" {
				if response.Body.String() != "OK" {
					t.Fatalf("self health = %q", response.Body.String())
				}
				return
			}
			var body struct {
				Service string `json:"service"`
				Status  string `json:"status"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Service != httptest.NewRequest(http.MethodGet, "/health"+tc.query, nil).URL.Query().Get("service") || body.Status != tc.status {
				t.Fatalf("unexpected response: %s", response.Body.String())
			}
		})
	}
}
