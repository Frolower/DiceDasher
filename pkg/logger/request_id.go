package logger

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const reqIDKey ctxKey = "req_id"

func NewReqID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return id
}

func WithReqID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, reqIDKey, id)
}

func ReqIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(reqIDKey)
	id, ok := v.(uuid.UUID)
	return id, ok
}
