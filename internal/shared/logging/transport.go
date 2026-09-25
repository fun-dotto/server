package logging

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const maxLoggedBodyBytes = 4096

// Transport は外部 API 呼び出しをログに出す http.RoundTripper。
// 呼び出し元の trace を X-Cloud-Trace-Context で上流へ伝播し、
// 上流が 4xx/5xx を返した場合はレスポンス本文もログに含める。
type Transport struct {
	Base http.RoundTripper
}

// WrapClient は client の Transport を Transport で包む。
func WrapClient(client *http.Client) *http.Client {
	base := client.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	client.Transport = &Transport{Base: base}
	return client
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if tr, ok := traceFrom(ctx); ok && req.Header.Get(cloudTraceHeader) == "" {
		req = req.Clone(ctx)
		v := tr.traceID
		if tr.spanID != "" {
			v += "/" + tr.spanID
		}
		if tr.sampled {
			v += ";o=1"
		}
		req.Header.Set(cloudTraceHeader, v)
	}

	start := time.Now()
	resp, err := t.Base.RoundTrip(req)
	latency := slog.String("latency", fmt.Sprintf("%.6fs", time.Since(start).Seconds()))
	target := req.Method + " " + req.URL.String()

	if err != nil {
		slog.ErrorContext(ctx, "upstream request failed: "+target, latency, slog.String("error", err.Error()))
		return nil, err
	}

	attrs := []any{latency, slog.Int("status", resp.StatusCode)}
	level := slog.LevelDebug
	if resp.StatusCode >= http.StatusBadRequest {
		level = slog.LevelWarn
		if resp.StatusCode >= http.StatusInternalServerError {
			level = slog.LevelError
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxLoggedBodyBytes))
		resp.Body = struct {
			io.Reader
			io.Closer
		}{io.MultiReader(bytes.NewReader(body), resp.Body), resp.Body}
		if readErr == nil {
			attrs = append(attrs, slog.String("responseBody", string(body)))
		}
	}
	slog.Log(ctx, level, fmt.Sprintf("upstream %s -> %d", target, resp.StatusCode), attrs...)
	return resp, nil
}
