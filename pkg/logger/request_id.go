package logger

import "context"

type ctxKey string

const reqIDKey ctxKey = "req_id"

func WithReqID(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, reqIDKey, id)
}

func ReqIDFromContext(ctx context.Context) (uint64, bool) {
	v := ctx.Value(reqIDKey)
	id, ok := v.(uint64)
	return id, ok
}
