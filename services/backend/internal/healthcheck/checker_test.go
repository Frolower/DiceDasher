package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPHealth(t *testing.T) {
	for _, code := range []int{200, 204, 302, 500, 503} {
		t.Run(http.StatusText(code), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/health" || r.Method != http.MethodGet {
					t.Error("unexpected health request")
				}
				w.Header().Set("Location", "/redirected")
				w.WriteHeader(code)
			}))
			defer server.Close()
			checker, err := New(map[string]string{"character": server.URL + "/health"})
			if err != nil {
				t.Fatal(err)
			}
			defer checker.Close()
			err = checker.Check(context.Background(), "character")
			if (err == nil) != (code == http.StatusOK) {
				t.Fatalf("Check for HTTP %d = %v", code, err)
			}
			if err := checker.Check(context.Background(), "missing"); !errors.Is(err, ErrUnknownService) {
				t.Fatalf("unknown service = %v", err)
			}
		})
	}
}

func TestHTTPHealthTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	checker, err := New(map[string]string{"slow": server.URL})
	if err != nil {
		t.Fatal(err)
	}
	defer checker.Close()
	started := time.Now()
	if err := checker.Check(context.Background(), "slow"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout = %v", err)
	}
	if elapsed := time.Since(started); elapsed > 4*time.Second {
		t.Fatalf("health timeout took %s", elapsed)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := checker.Check(ctx, "slow"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel = %v", err)
	}
}

func TestInvalidTargets(t *testing.T) {
	for _, target := range []map[string]string{
		{"backend": "http://localhost/health"},
		{"": "http://localhost/health"},
		{"room": "localhost:50051"},
		{"room": "ftp://localhost"},
		{"room": "grpc://localhost:50051/path"},
		{"room": "http://user:password@localhost/health"},
	} {
		if checker, err := New(target); err == nil {
			checker.Close()
			t.Errorf("accepted invalid target %v", target)
		}
	}
}

func TestHTTPHealthBodyTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()
	checker, err := New(map[string]string{"slow": server.URL})
	if err != nil {
		t.Fatal(err)
	}
	defer checker.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := checker.Check(ctx, "slow"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("body timeout = %v", err)
	}
}
