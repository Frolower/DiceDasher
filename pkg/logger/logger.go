package logger

import (
	"context"
	"log"
	"os"
)

const (
	cReset = "\033[0m"

	// foreground
	cGray         = "\033[90m"
	cRed          = "\033[31m"
	cGreen        = "\033[32m"
	cBlue         = "\033[34m"
	cBrightYellow = "\033[93m"
)

func Logf(ctx context.Context, format string, args ...any) {
	logf(ctx, "LOG", cGray, format, args...)
}

func logf(ctx context.Context, label, color, format string, args ...any) {
	if id, ok := ReqIDFromContext(ctx); ok {
		format = "id=%s | " + format
		args = append([]any{id.String()}, args...)
	}
	separator := "| "
	if label == "REQUEST" || label == "RESPONSE" {
		separator = "|"
	}
	log.Printf("|%s%-8s%s"+separator+format, append([]any{color, label, cReset}, args...)...)
}

func LogWarning(message string) {
	LogWarningf(context.Background(), "%s", message)
}

func LogWarningf(ctx context.Context, format string, args ...any) {
	logf(ctx, "WARNING", cBrightYellow, format, args...)
}

func LogErrorf(ctx context.Context, format string, args ...any) {
	logf(ctx, "ERROR", cRed, format, args...)
}

func LogRequestf(ctx context.Context, format string, args ...any) {
	logf(ctx, "REQUEST", cGreen, format, args...)
}

func LogResponsef(ctx context.Context, format string, args ...any) {
	logf(ctx, "RESPONSE", cBlue, format, args...)
}

func LogFatalf(ctx context.Context, format string, args ...any) {
	LogErrorf(ctx, format, args...)
	os.Exit(1)
}
