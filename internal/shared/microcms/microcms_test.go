package microcms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// fakeServer は total 件の記事を offset / limit に従って返す。
func fakeServer(t *testing.T, total int, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-MICROCMS-API-KEY") != "key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		contents := []map[string]any{}
		for i := offset; i < total && i < offset+limit; i++ {
			contents = append(contents, map[string]any{
				"id":          fmt.Sprintf("a%d", i),
				"title":       fmt.Sprintf("title %d", i),
				"publishedAt": time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
			})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"contents": contents, "totalCount": total, "offset": offset, "limit": limit})
	}))
}

func TestListArticles(t *testing.T) {
	tests := []struct {
		name    string
		total   int
		status  int
		want    int
		wantErr bool
	}{
		{name: "0 件", total: 0, status: http.StatusOK, want: 0},
		{name: "1 ページに収まる", total: 100, status: http.StatusOK, want: 100},
		{name: "複数ページにまたがる", total: 101, status: http.StatusOK, want: 101},
		{name: "2xx 以外はエラー", total: 1, status: http.StatusInternalServerError, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := fakeServer(t, tt.total, tt.status)
			defer srv.Close()
			got, err := newClient(srv.URL, "key", srv.Client()).ListArticles(context.Background(), PublishedAtDesc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Errorf("len = %d, want %d", len(got), tt.want)
			}
		})
	}
}
