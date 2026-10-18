package logger

import (
	"context"
	"log"
	"os"
)

const (
	cReset = "\033[0m"

	// foreground
	cBlack        = "\033[30m"
	cWhite        = "\033[97m"
	cGray         = "\033[90m"
	cRed          = "\033[31m"
	cBrightYellow = "\033[93m"

	bgBrightRed = "\033[101m"
)

func Logf(ctx context.Context, format string, args ...any) {
	if id, ok := ReqIDFromContext(ctx); ok {
		log.Printf("| id=%s | "+format, append([]any{id.String()}, args...)...)
		return
	}
	log.Printf(format, args...)
}

func LogWarning(message string) {
	log.Printf("%sWARNING: %s%s", cBrightYellow, message, cReset)
}

func LogWarningf(ctx context.Context, format string, args ...any) {
	Logf(ctx, cBrightYellow+"WARNING: "+format+cReset, args...)
}

func LogErrorf(ctx context.Context, format string, args ...any) {
	Logf(ctx, cRed+"ERROR: "+format+cReset, args...)
}

func LogFatalf(ctx context.Context, format string, args ...any) {
	Logf(ctx, bgBrightRed+cWhite+"FATAL: "+format+cReset, args...)
	os.Exit(1)
}
