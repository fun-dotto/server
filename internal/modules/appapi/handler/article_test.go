package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/academiccalendar"
	"github.com/fun-dotto/server/internal/shared/microcms"
)

type fakeArticles struct {
	articles []microcms.Article
	err      error
}

func (f fakeArticles) ListArticles(context.Context, microcms.SortOrder) ([]microcms.Article, error) {
	return f.articles, f.err
}

func TestListArticles(t *testing.T) {
	published := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		articles ArticleLister
		wantCode connect.Code
		want     *appv1.Article
	}{
		{
			name:     "記事の ID・URL・タイトル・公開日時を返す",
			articles: fakeArticles{articles: []microcms.Article{{ID: "abc", Title: "t", PublishedAt: published}}},
			want:     &appv1.Article{Id: "abc", Url: "https://fun-dotto.com/articles/abc", Title: "t", PublishTime: toTimestamp(published)},
		},
		{name: "microCMS の障害は UNAVAILABLE", articles: fakeArticles{err: errors.New("boom")}, wantCode: connect.CodeUnavailable},
		{name: "microCMS が未設定なら UNAVAILABLE", articles: nil, wantCode: connect.CodeUnavailable},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(nil, academiccalendar.Calendar{}, tt.articles, "https://fun-dotto.com")
			res, err := h.ListArticles(context.Background(), connect.NewRequest(&appv1.ListArticlesRequest{}))
			if tt.wantCode != 0 {
				wantCode(t, err, tt.wantCode)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			got := res.Msg.GetArticles()
			if len(got) != 1 || got[0].GetId() != tt.want.GetId() || got[0].GetUrl() != tt.want.GetUrl() ||
				got[0].GetTitle() != tt.want.GetTitle() || !got[0].GetPublishTime().AsTime().Equal(published) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
