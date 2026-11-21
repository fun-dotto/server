package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/shared/microcms"
	"github.com/fun-dotto/server/internal/shared/uuidv5"
)

// ArticleLister は公開中の記事を取得する。
type ArticleLister interface {
	ListArticles(ctx context.Context, order microcms.SortOrder) ([]microcms.Article, error)
}

// AnnouncementRepository は microCMS の記事をおしらせとして返す。
// アプリの後方互換性のため、OpenAPI のスキーマとドメインモデルは記事導入前のままにしている。
type AnnouncementRepository struct {
	articles       ArticleLister
	articleBaseURL string
}

// NewAnnouncementRepository は AnnouncementRepository を返す。articleBaseURL は記事ページの URL のベース。
func NewAnnouncementRepository(articles ArticleLister, articleBaseURL string) *AnnouncementRepository {
	return &AnnouncementRepository{articles: articles, articleBaseURL: articleBaseURL}
}

// GetAnnouncements は公開中の記事を返す。
// microCMS は公開中の記事しか返さないため、query.FilterIsActive は常に満たされる。
func (r *AnnouncementRepository) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	order := microcms.PublishedAtDesc
	if query.SortByDate != nil && *query.SortByDate == domain.SortDirectionAsc {
		order = microcms.PublishedAtAsc
	}
	articles, err := r.articles.ListArticles(ctx, order)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Announcement, len(articles))
	for i, a := range articles {
		result[i] = toDomainAnnouncement(a, r.articleBaseURL)
	}
	return result, nil
}

// toDomainAnnouncement は記事をおしらせに変換する。
// アプリは id を UUID として扱うため、記事 URL から決定的な UUID を生成する。
func toDomainAnnouncement(a microcms.Article, baseURL string) domain.Announcement {
	url := microcms.ArticleURL(baseURL, a.ID)
	return domain.Announcement{
		ID:            uuidv5.NewSHA1(uuidv5.NamespaceURL, []byte(url)).String(),
		Title:         a.Title,
		AvailableFrom: a.PublishedAt,
		URL:           url,
	}
}
