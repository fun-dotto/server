// Package logging はサーバー全体で共通の構造化ログ (log/slog) を提供する。
// Cloud Run 上では Cloud Logging が解釈できる JSON (severity / message /
// logging.googleapis.com/trace) を標準出力に書き、同一リクエストのログを
// trace で束ねる。ローカルでは読みやすいテキスト形式で出力する。
package logging

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
)

const (
	traceKey   = "logging.googleapis.com/trace"
	spanKey    = "logging.googleapis.com/spanId"
	sampledKey = "logging.googleapis.com/trace_sampled"
)

var projectID string

// Setup は slog の既定ロガーを初期化する。標準の log パッケージの出力も
// slog 経由になるため、既存の log.Printf も同じ形式で出力される。
func Setup() {
	var h slog.Handler
	if onCloudRun() {
		projectID = resolveProjectID()
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelDebug,
			ReplaceAttr: replaceCloudLoggingAttr,
		})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(&traceHandler{Handler: h}))
	log.SetFlags(0)
}

func onCloudRun() bool {
	return os.Getenv("K_SERVICE") != "" || os.Getenv("CLOUD_RUN_JOB") != ""
}

func resolveProjectID() string {
	if id := os.Getenv("GOOGLE_CLOUD_PROJECT"); id != "" {
		return id
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	id, err := metadata.ProjectIDWithContext(ctx)
	if err != nil {
		return ""
	}
	return id
}

// replaceCloudLoggingAttr は slog の標準キーを Cloud Logging の特殊フィールドに変換する。
func replaceCloudLoggingAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return a
	}
	switch a.Key {
	case slog.MessageKey:
		a.Key = "message"
	case slog.LevelKey:
		a.Key = "severity"
		if lv, ok := a.Value.Any().(slog.Level); ok && lv == slog.LevelWarn {
			a.Value = slog.StringValue("WARNING")
		}
	case slog.SourceKey:
		a.Key = "logging.googleapis.com/sourceLocation"
	}
	return a
}

type traceInfo struct {
	traceID string
	spanID  string
	sampled bool
}

type traceCtxKey struct{}

func withTrace(ctx context.Context, t traceInfo) context.Context {
	return context.WithValue(ctx, traceCtxKey{}, t)
}

func traceFrom(ctx context.Context) (traceInfo, bool) {
	if ctx == nil {
		return traceInfo{}, false
	}
	t, ok := ctx.Value(traceCtxKey{}).(traceInfo)
	return t, ok && t.traceID != ""
}

// traceHandler は context に載った trace 情報をログに付与する。
type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if t, ok := traceFrom(ctx); ok {
		trace := t.traceID
		if projectID != "" {
			trace = "projects/" + projectID + "/traces/" + t.traceID
		}
		r.AddAttrs(slog.String(traceKey, trace))
		if t.spanID != "" {
			r.AddAttrs(slog.String(spanKey, t.spanID))
		}
		r.AddAttrs(slog.Bool(sampledKey, t.sampled))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithGroup(name)}
}
