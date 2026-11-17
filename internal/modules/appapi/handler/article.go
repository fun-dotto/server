package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/microcms"
)

// ListArticles は公開中の記事の一覧を、公開日時の新しい順に返す。
func (h *Handler) ListArticles(ctx context.Context, _ *connect.Request[appv1.ListArticlesRequest]) (*connect.Response[appv1.ListArticlesResponse], error) {
	articles, err := h.articles.ListArticles(ctx, microcms.PublishedAtDesc)
	if err != nil {
		return nil, internalError(ctx, "ListArticles", err)
	}
	return connect.NewResponse(&appv1.ListArticlesResponse{
		Articles: mapSlice(articles, func(a microcms.Article) *appv1.Article {
			return &appv1.Article{
				Id:          a.ID,
				Url:         microcms.ArticleURL(h.articleBaseURL, a.ID),
				Title:       a.Title,
				PublishTime: toTimestamp(a.PublishedAt),
			}
		}),
	}), nil
}
