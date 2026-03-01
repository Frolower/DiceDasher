package system

import (
	"context"
	"encoding/json"
)

type Resolver interface {
	Resolve(ctx context.Context, raw json.RawMessage) (any, int, error)
}
