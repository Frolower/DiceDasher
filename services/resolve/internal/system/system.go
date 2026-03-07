package system

import (
	"context"
	"encoding/json"
)

type Resolver interface {
	Resolve(ctx context.Context, action string, raw json.RawMessage) (any, int, error)
}
