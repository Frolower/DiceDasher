package logger

import (
	"bytes"
	"context"
	"log"
	"testing"
)

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
			want := tc.color + tc.label + padding + cReset + "| id=" + id.String() + " | value=42\n"
			if got := buf.String(); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}
