package httputil

import (
	"encoding/json/jsontext"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUnpackJSON(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		wantErr    bool
	}{
		{"valid", `{"name":"Mira"}`, false},
		{"empty", ``, true},
		{"unknown field", `{"extra":1}`, true},
		{"duplicate name", `{"name":"Mira","name":"Other"}`, true},
		{"invalid UTF8", "{\"name\":\"\xff\"}", true},
		{"multiple values", `{"name":"Mira"} {}`, true},
		{"trailing junk", `{"name":"Mira"}x`, true},
		{"trailing whitespace", "{\"name\":\"Mira\"} \n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var dst struct {
				Name string `json:"name"`
			}
			err := UnpackJSON(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body)), &dst)
			if (err != nil) != tc.wantErr {
				t.Fatalf("UnpackJSON error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUnpackRawJSON(t *testing.T) {
	for _, body := range []string{"", " \n", `{"nested":{"value":1}}`} {
		var dst jsontext.Value
		if err := UnpackJSON(httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body)), &dst); err != nil {
			t.Fatalf("body %q: %v", body, err)
		}
		if strings.TrimSpace(body) != "" && string(dst) != body {
			t.Fatalf("got %s, want %s", dst, body)
		}
	}
}

func TestPackJSON(t *testing.T) {
	w := httptest.NewRecorder()
	if err := PackJSON(w, http.StatusCreated, struct {
		Payload jsontext.Value `json:"payload"`
		Text    string         `json:"text"`
	}{jsontext.Value(`{"ok":true}`), "<>&"}); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusCreated || w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected response: %v", w.Result())
	}
	if got, want := w.Body.String(), "{\"payload\":{\"ok\":true},\"text\":\"\\u003c\\u003e\\u0026\"}\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
