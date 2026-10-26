package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMultipleOrigins(t *testing.T) {
	for _, origin := range []string{"http://localhost:8888", "https://second.example", "http://localhost:8000", "https://second.example.evil", "null", ""} {
		for _, method := range []string{http.MethodGet, http.MethodOptions} {
			t.Run(method+origin, func(t *testing.T) {
				called := false
				h := CORS([]string{"http://localhost:8888", "https://second.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
				r := httptest.NewRequest(method, "/health", nil)
				r.Header.Set("Origin", origin)
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				want := ""
				if origin == "http://localhost:8888" || origin == "https://second.example" {
					want = origin
				}
				if got := w.Header().Get("Access-Control-Allow-Origin"); got != want {
					t.Fatalf("allow origin = %q, want %q", got, want)
				}
				credentials := w.Header().Get("Access-Control-Allow-Credentials")
				if (want != "" && credentials != "true") || (want == "" && credentials != "") {
					t.Fatalf("unexpected credentials: %q", credentials)
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
