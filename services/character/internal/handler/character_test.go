package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"diceDasher/services/character/internal/system"
	"diceDasher/services/character/internal/system/tes"
)

func TestCharacterValidationResponse(t *testing.T) {
	system.Register("tes", tes.Character{})
	req := httptest.NewRequest(http.MethodPost, "/character?system=tes", strings.NewReader(`{"user_id":"d7a92c4c-7c65-41d8-946a-4b94d3e721f9","type":"pc","character":{"name":""}}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	postUserCreatedCharacterHandler(w, req)
	if w.Code != 422 {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Fatal("expected JSON")
	}
	var response system.ValidationError
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Violations) != 1 || response.Violations[0].Field != "character.name" || response.Violations[0].Code != "required" {
		t.Fatalf("unexpected violations: %+v", response)
	}
}
