package httputil

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const (
	cReset = "\033[0m"

	// foreground
	cBlack = "\033[30m"
	cWhite = "\033[97m"
	cGray  = "\033[90m"

	// background
	bgGreen  = "\033[42m"
	bgCyan   = "\033[46m"
	bgBlue   = "\033[44m"
	bgYellow = "\033[43m"
	bgRed    = "\033[41m"
	bgGray   = "\033[100m"

	// background bright
	bgBrightBlack   = "\033[100m"
	bgBrightRed     = "\033[101m"
	bgBrightGreen   = "\033[102m"
	bgBrightYellow  = "\033[103m"
	bgBrightBlue    = "\033[104m"
	bgBrightMagenta = "\033[105m"
	bgBrightCyan    = "\033[106m"
	bgBrightWhite   = "\033[107m"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

type ctxKey string

const reqIDKey ctxKey = "req_id"

var reqCounter uint64

func nextReqID() uint64 {
	return atomic.AddUint64(&reqCounter, 1)
}

func ReqIDFromContext(ctx context.Context) (uint64, bool) {
	v := ctx.Value(reqIDKey)
	id, ok := v.(uint64)
	return id, ok
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	// If WriteHeader wasn't called, default is 200
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func styleForStatus(code int, mode string) string {
	bg200, bg300, bg400, bg500 := bgGreen, bgCyan, bgYellow, bgRed

	switch mode {
	case "contrast":
		bg200, bg300, bg400, bg500 = bgBrightGreen, bgBrightCyan, bgBrightYellow, bgBrightRed
	case "default", "":
		// keep defaults
	default:
		// unknown mode -> default palette
	}

	switch {
	case code >= 200 && code < 300:
		return bg200 + cWhite
	case code >= 300 && code < 400:
		return bg300 + cWhite
	case code >= 400 && code < 500:
		return bg400 + cWhite
	default:
		return bg500 + cWhite
	}
}

func styleForMethod(method, mode string) string {
	// Pick palette based on mode
	getBg, postBg, putBg, delBg, defBg := bgBlue, bgCyan, bgYellow, bgRed, bgGray

	switch mode {
	case "contrast":
		getBg, postBg, putBg, delBg = bgBrightBlue, bgBrightCyan, bgBrightYellow, bgBrightRed
	case "default", "":
		// keep defaults
	default:
		// unknown mode -> default palette
	}

	switch strings.ToUpper(method) {
	case http.MethodGet:
		return getBg + cWhite
	case http.MethodPost:
		return postBg + cWhite
	case http.MethodPut:
		return putBg + cWhite
	case http.MethodDelete:
		return delBg + cWhite
	default:
		return defBg + cWhite
	}
}

func RequestLogger(next http.Handler) http.Handler {
	return RequestLoggerWithMode(next, "normal")
}

func RequestLoggerWithMode(next http.Handler, mode string) http.Handler {
	switch mode {
	case "", "default", "contrast":
		if mode == "" {
			mode = "default"
		}
	default:
		mode = "default"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := nextReqID()

		// Put id into context so handlers can use it if they want
		r = r.WithContext(context.WithValue(r.Context(), reqIDKey, id))

		sw := &statusWriter{ResponseWriter: w}
		start := time.Now()

		ms := styleForMethod(r.Method, mode)
		ip := clientIP(r)
		ct := r.Header.Get("Content-Type")
		if ct == "" {
			ct = "none"
		}
		cl := r.ContentLength // can be -1

		// Log on request
		log.Printf("| %-39s | ct=%-32s cl=%-10d | id=%06d |%s%-6s%s| \"%s\"",
			ip, ct, cl, id, ms, r.Method,
			cReset, r.URL.RequestURI())

		next.ServeHTTP(sw, r)

		dur := time.Since(start)
		if sw.status == 0 {
			sw.status = http.StatusOK
		}

		ss := styleForStatus(sw.status, mode)

		// Log on response
		log.Printf("|%s %3d %s| %22s | %60s | id=%06d |%s%-6s%s| \"%s\"%s bytes=%d",
			ss, sw.status, cReset,
			dur, r.UserAgent(), id,
			ms, r.Method, cReset,
			r.URL.RequestURI(), cGray,
			sw.bytes)
	})
}
