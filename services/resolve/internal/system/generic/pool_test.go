package generic

import (
	"context"
	"diceDasher/pkg/dice"
	"diceDasher/services/resolve/internal/system"
	"encoding/json"
	"testing"
)

func TestRollPoolValidation(t *testing.T) {
	for _, req := range []request{{Number: -1, Size: 6}, {Number: dice.MaxDice + 1, Size: 6}, {Number: 1, Size: dice.MaxSides + 1}} {
		raw, _ := json.Marshal(req)
		_, err := (Resolver{}).Resolve(context.Background(), "roll", raw)
		if !system.IsValidation(err) {
			t.Fatalf("accepted %+v", req)
		}
	}
}
