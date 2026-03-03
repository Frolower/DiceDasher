package logger

import (
	"context"
	"log"
)

const (
	cReset = "\033[0m"

	// foreground
	cBlack        = "\033[30m"
	cWhite        = "\033[97m"
	cGray         = "\033[90m"
	cBrightYellow = "\033[93m"
)

func Logf(ctx context.Context, format string, args ...any) {
	if id, ok := ReqIDFromContext(ctx); ok {
		log.Printf("id=%s "+format, append([]any{id.String()}, args...)...)
		return
	}
	log.Printf(format, args...)
}

func LogWarning(message string) {
	log.Printf("%sWARNING: %s%s", cBrightYellow, message, cReset)
}

func LogWarningf(format string, args ...any) {
	log.Printf(cBrightYellow+"WARNING: "+format+cReset, args...)
}
