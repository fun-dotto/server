package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const cloudTraceHeader = "X-Cloud-Trace-Context"

// Middleware は http.Handler を包み、リクエストごとに trace を context に載せ、
// 完了時にアクセスログを出す。ステータス 5xx は ERROR、4xx は WARNING、
// それ以外は INFO で出力する。RecordError で記録されたエラーも併せて出力し、
// パニックはスタックトレース付きでログに出して 500 を返す。
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()
		if t, ok := parseTrace(r.Header); ok {
			ctx = withTrace(ctx, t)
		}
		errs := &requestErrors{}
		ctx = context.WithValue(ctx, requestErrorsKey{}, errs)
		r = r.WithContext(ctx)
		rw := &statusRecorder{ResponseWriter: w}

		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				slog.ErrorContext(ctx, fmt.Sprintf("panic recovered: %v", rec),
					slog.String("stack", string(debug.Stack())),
				)
				errs.add(fmt.Errorf("panic: %v", rec))
				if !rw.wroteHeader {
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
			logAccess(ctx, r, rw, time.Since(start), errs.join())
		}()

		next.ServeHTTP(rw, r)
	})
}

func logAccess(ctx context.Context, r *http.Request, rw *statusRecorder, latency time.Duration, err error) {
	status := rw.status()
	level := slog.LevelInfo
	switch {
	case status >= http.StatusInternalServerError:
		level = slog.LevelError
	case status >= http.StatusBadRequest:
		level = slog.LevelWarn
	}

	attrs := []slog.Attr{
		slog.Group("httpRequest",
			slog.String("requestMethod", r.Method),
			slog.String("requestUrl", r.URL.RequestURI()),
			slog.Int("status", status),
			slog.Int("responseSize", rw.size),
			slog.String("userAgent", r.UserAgent()),
			slog.String("remoteIp", r.RemoteAddr),
			slog.String("latency", fmt.Sprintf("%.6fs", latency.Seconds())),
		),
	}
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	msg := fmt.Sprintf("%s %s -> %d", r.Method, r.URL.Path, status)
	slog.LogAttrs(ctx, level, msg, attrs...)
}

// RecordError はリクエスト処理中のエラーを記録し、アクセスログに含める。
// ctx は Middleware を通ったリクエストの context である必要がある。
func RecordError(ctx context.Context, err error) {
	if err == nil || ctx == nil {
		return
	}
	if errs, ok := ctx.Value(requestErrorsKey{}).(*requestErrors); ok {
		errs.add(err)
		return
	}
	slog.ErrorContext(ctx, err.Error())
}

type requestErrorsKey struct{}

type requestErrors struct {
	mu   sync.Mutex
	errs []error
}

func (e *requestErrors) add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errs = append(e.errs, err)
}

func (e *requestErrors) join() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return errors.Join(e.errs...)
}

// statusRecorder はレスポンスのステータスとサイズを記録する。
type statusRecorder struct {
	http.ResponseWriter
	code        int
	size        int
	wroteHeader bool
}

func (w *statusRecorder) WriteHeader(code int) {
	if !w.wroteHeader {
		w.code = code
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *statusRecorder) status() int {
	if !w.wroteHeader {
		return http.StatusOK
	}
	return w.code
}

// Unwrap は http.ResponseController が元の ResponseWriter の機能
// (Flush など) を使えるようにする。
func (w *statusRecorder) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// parseTrace は X-Cloud-Trace-Context ("TRACE_ID/SPAN_ID;o=1") または
// W3C traceparent ("00-TRACE_ID-SPAN_ID-FLAGS") から trace 情報を取り出す。
func parseTrace(h http.Header) (traceInfo, bool) {
	if v := h.Get(cloudTraceHeader); v != "" {
		var t traceInfo
		rest := v
		if i := strings.Index(rest, ";"); i >= 0 {
			t.sampled = strings.Contains(rest[i:], "o=1")
			rest = rest[:i]
		}
		t.traceID, t.spanID, _ = strings.Cut(rest, "/")
		return t, t.traceID != ""
	}
	if v := h.Get("traceparent"); v != "" {
		parts := strings.Split(v, "-")
		if len(parts) == 4 {
			return traceInfo{traceID: parts[1], spanID: parts[2], sampled: strings.HasSuffix(parts[3], "1")}, true
		}
	}
	return traceInfo{}, false
}
