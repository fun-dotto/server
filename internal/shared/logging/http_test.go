package logging

import (
	"net/http"
	"testing"
)

func TestParseTrace(t *testing.T) {
	tests := []struct {
		name   string
		header http.Header
		want   traceInfo
		wantOK bool
	}{
		{
			name:   "X-Cloud-Trace-Context (sampled)",
			header: http.Header{"X-Cloud-Trace-Context": {"abc123/456;o=1"}},
			want:   traceInfo{traceID: "abc123", spanID: "456", sampled: true},
			wantOK: true,
		},
		{
			name:   "X-Cloud-Trace-Context (trace のみ)",
			header: http.Header{"X-Cloud-Trace-Context": {"abc123"}},
			want:   traceInfo{traceID: "abc123"},
			wantOK: true,
		},
		{
			name:   "traceparent",
			header: http.Header{"Traceparent": {"00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}},
			want:   traceInfo{traceID: "0af7651916cd43dd8448eb211c80319c", spanID: "b7ad6b7169203331", sampled: true},
			wantOK: true,
		},
		{
			name:   "ヘッダなし",
			header: http.Header{},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseTrace(tt.header)
			if ok != tt.wantOK || got != tt.want {
				t.Errorf("parseTrace() = %+v, %v; want %+v, %v", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
