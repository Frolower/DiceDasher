package httphelper_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httphelper "github.com/frolower/dicedasher/shared/http"
)

type testStruct struct {
	Number  int  `json:"number"`
	Boolean bool `json:"bool"`
}

func TestWriteJSON(t *testing.T) {
	input := testStruct{Number: 6, Boolean: true}
	expected := `{
  "number": 6,
  "bool": true
}`

	body, status, err := httphelper.WriteJSON(http.StatusOK, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

	if string(body) != expected {
		t.Errorf("unexpected JSON output:\nexpected:\n%s\ngot:\n%s", expected, string(body))
	}
}

func TestReadJSON_Valid(t *testing.T) {
	body := `{"number":42, "bool":true}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	result, err := httphelper.ReadJSON[testStruct](req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Number != 42 || result.Boolean != true {
		t.Errorf("unexpected decoded result: %+v", result)
	}
}

func TestReadJSON_MalformedJSON(t *testing.T) {
	body := `{"number":42, "bool":` // malformed
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	_, err := httphelper.ReadJSON[testStruct](req)
	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
}

func TestReadJSON_UnknownField(t *testing.T) {
	body := `{"number":42, "bool":true, "extra":123}` // unknown field
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	_, err := httphelper.ReadJSON[testStruct](req)
	if err == nil {
		t.Error("expected error for unknown field, got nil")
	}
}

func TestReadJSON_WrongType(t *testing.T) {
	body := `{"number":"forty-two", "bool":true}` // "number" should be int
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	_, err := httphelper.ReadJSON[testStruct](req)
	if err == nil {
		t.Error("expected error for type mismatch, got nil")
	}
}

func TestReadJSON_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	_, err := httphelper.ReadJSON[testStruct](req)
	if err == nil {
		t.Error("expected error for empty body, got nil")
	}
}
