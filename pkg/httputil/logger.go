package httputil

import (
	"context"
	"log"
)

func Logf(ctx context.Context, format string, args ...any) {
	if id, ok := ReqIDFromContext(ctx); ok {
		// prefix the format, keep args untouched
		log.Printf("id=%06d "+format, append([]any{id}, args...)...)
		return
	}
	log.Printf(format, args...)
}
