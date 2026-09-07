package generic

import (
	"context"
	"diceDasher/pkg/dice"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRollPoolValidation(t *testing.T) {
	for _, req := range []request{{Number: -1, Size: 6}, {Number: dice.MaxDice + 1, Size: 6}, {Number: 1, Size: dice.MaxSides + 1}} {
		raw, _ := json.Marshal(req)
		_, status, err := (Resolver{}).Resolve(context.Background(), "roll", raw)
		if err == nil || status != http.StatusUnprocessableEntity {
			t.Fatalf("accepted %+v", req)
		}
	}
}
