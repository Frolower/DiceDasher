package httputil

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseAllowedOrigins(t *testing.T) {
	for _, tc := range []struct {
		name, raw string
		want      []string
	}{
		{"default", "", []string{"http://localhost:8888", "http://localhost:8000"}},
		{"blank", "  ", []string{"http://localhost:8888", "http://localhost:8000"}},
		{"trimmed", " http://localhost:8888, https://example.com ", []string{"http://localhost:8888", "https://example.com"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseAllowedOrigins(tc.raw, "http://localhost:8888,http://localhost:8000")
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	for _, raw := range []string{"*", "null", "localhost:8888", "ftp://example.com", "https://user@example.com", "https://example.com/", "https://example.com?q=1", "https://example.com?", "https://example.com#", "https://*.example.com", "http://:8000", "http://localhost:8888,", ",http://localhost:8888"} {
		t.Run(raw, func(t *testing.T) {
			if _, err := ParseAllowedOrigins(raw, "http://localhost:8888"); err == nil {
				t.Fatal("accepted invalid origin")
			}
		})
	}
}

func TestCORSMultipleOrigins(t *testing.T) {
	for _, origin := range []string{"http://localhost:8888", "http://localhost:8000", "https://evil.example", "http://localhost:8000.evil.example", "null", ""} {
		for _, method := range []string{http.MethodGet, http.MethodOptions} {
			t.Run(method+origin, func(t *testing.T) {
				called := false
				h := CORS([]string{"http://localhost:8888", "http://localhost:8000"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
				r := httptest.NewRequest(method, "/health", nil)
				r.Header.Set("Origin", origin)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				want := ""
				if origin == "http://localhost:8888" || origin == "http://localhost:8000" {
					want = origin
				}
				if got := w.Header().Get("Access-Control-Allow-Origin"); got != want {
					t.Fatalf("allow origin = %q, want %q", got, want)
				}
				if w.Header().Get("Vary") != "Origin" {
					t.Fatal("missing Vary: Origin")
				}
				if method == http.MethodOptions {
					if called || w.Code != http.StatusNoContent {
						t.Fatal("preflight reached handler or returned wrong status")
					}
				} else if !called {
					t.Fatal("GET did not reach handler")
				}
			})
		}
	}
}
