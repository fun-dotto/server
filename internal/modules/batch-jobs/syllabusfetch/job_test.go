package syllabusfetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobFetchDetailWithRetry(t *testing.T) {
	tests := []struct {
		name      string
		statuses  []int
		wantErr   bool
		wantCalls int
	}{
		{name: "5xx は 1 回やり直して成功する", statuses: []int{http.StatusBadGateway, http.StatusOK}, wantCalls: 2},
		{name: "5xx が続けばエラー", statuses: []int{http.StatusBadGateway, http.StatusBadGateway}, wantErr: true, wantCalls: 2},
		{name: "4xx はやり直さない", statuses: []int{http.StatusNotFound}, wantErr: true, wantCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statuses[calls])
				calls++
			}))
			t.Cleanup(srv.Close)
			job := &Job{client: NewClient(srv.URL, 0), retryWait: 0}
			_, err := job.fetchDetailWithRetry(context.Background(), "100101", 2026)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}
