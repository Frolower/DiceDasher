package logger

import (
	"bytes"
	"context"
	"log"
	"testing"
	"time"
)

// TestLogPreview prints colored examples when run with go test -v -run TestLogPreview.
func TestLogPreview(t *testing.T) {
	ctx := context.Background()
	id := NewReqID()
	Logf(ctx, "backend service READY on :8080")
	LogWarningf(ctx, "example warning: slow database connection")
	LogErrorf(ctx, "example error: database unavailable")
	LogRequestf(ctx, "\033[46m\033[97m%-6s%s| %-10s | id=%s | %s",
		"POST", cReset, "/login", id.String(), "127.0.0.1")
	LogResponsef(ctx, "\033[42m\033[97m %3d %s| %-10s | id=%s",
		200, cReset, 69841*time.Microsecond, id.String())
}

func TestLogLabels(t *testing.T) {
	output, flags, prefix := log.Writer(), log.Flags(), log.Prefix()
	t.Cleanup(func() {
		log.SetOutput(output)
		log.SetFlags(flags)
		log.SetPrefix(prefix)
	})
	log.SetFlags(0)
	log.SetPrefix("")

	for _, tc := range []struct {
		label, color string
		write        func(context.Context, string, ...any)
	}{
		{"LOG", cGray, Logf},
		{"WARNING", cBrightYellow, LogWarningf},
		{"ERROR", cRed, LogErrorf},
		{"REQUEST", cGreen, LogRequestf},
		{"RESPONSE", cBlue, LogResponsef},
	} {
		t.Run(tc.label, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			ctx := WithReqID(context.Background(), NewReqID())
			id, _ := ReqIDFromContext(ctx)
			tc.write(ctx, "value=%d", 42)
			padding := ""
			for i := len(tc.label); i < 8; i++ {
				padding += " "
			}
			separator := "| "
			if tc.label == "REQUEST" || tc.label == "RESPONSE" {
				separator = "|"
			}
			want := "|" + tc.color + tc.label + padding + cReset + separator + "id=" + id.String() + " | value=42\n"
			if got := buf.String(); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}
