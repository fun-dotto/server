// Package microcms は microCMS のコンテンツ API から記事を取得するクライアント。
// 記事の入稿は microCMS の管理画面で行うため、読み取りのみを提供する。
package microcms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	httpClientTimeout = 10 * time.Second
	// microCMS のリスト取得で指定できる limit の上限。
	maxLimit = 100
)

// SortOrder は記事の並び順。
type SortOrder int

const (
	// PublishedAtDesc は公開日時の新しい順。
	PublishedAtDesc SortOrder = iota
	// PublishedAtAsc は公開日時の古い順。
	PublishedAtAsc
)

// Article は公開中の記事。本文はアプリでは使わないため取得しない。
type Article struct {
	ID          string
	Title       string
	PublishedAt time.Time
}

// Client は microCMS のコンテンツ API のクライアント。
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient は Client を返す。serviceDomain は https://<domain>.microcms.io の <domain>。
func NewClient(serviceDomain, apiKey string, httpClient *http.Client) *Client {
	return newClient("https://"+serviceDomain+".microcms.io", apiKey, httpClient)
}

func newClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: httpClientTimeout}
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), apiKey: apiKey, httpClient: httpClient}
}

// NewClientFromEnv は環境変数 MICROCMS_SERVICE_DOMAIN と MICROCMS_API_KEY から Client を返す。
func NewClientFromEnv() (*Client, error) {
	domain := os.Getenv("MICROCMS_SERVICE_DOMAIN")
	if domain == "" {
		return nil, errors.New("MICROCMS_SERVICE_DOMAIN is required")
	}
	apiKey := os.Getenv("MICROCMS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("MICROCMS_API_KEY is required")
	}
	return NewClient(domain, apiKey, nil), nil
}

type listArticlesResponse struct {
	Contents []struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		PublishedAt time.Time `json:"publishedAt"`
	} `json:"contents"`
	TotalCount int `json:"totalCount"`
}

// ListArticles は公開中の記事を全件返す。下書き・非公開の記事は microCMS が返さない。
func (c *Client) ListArticles(ctx context.Context, order SortOrder) ([]Article, error) {
	// publishedAt が同じ記事があっても offset ページングの境目で順序が揺れないよう、id を第 2 キーにする。
	orders := "-publishedAt,id"
	if order == PublishedAtAsc {
		orders = "publishedAt,id"
	}
	var articles []Article
	for {
		q := url.Values{}
		q.Set("fields", "id,title,publishedAt")
		q.Set("orders", orders)
		q.Set("limit", strconv.Itoa(maxLimit))
		q.Set("offset", strconv.Itoa(len(articles)))
		var res listArticlesResponse
		if err := c.get(ctx, "/api/v1/articles?"+q.Encode(), &res); err != nil {
			return nil, fmt.Errorf("list articles: %w", err)
		}
		for _, a := range res.Contents {
			articles = append(articles, Article{ID: a.ID, Title: a.Title, PublishedAt: a.PublishedAt})
		}
		if len(res.Contents) == 0 || len(articles) >= res.TotalCount {
			return articles, nil
		}
	}
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-MICROCMS-API-KEY", c.apiKey)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// ArticleURL は Web で記事を表示するページの URL を返す。baseURL は https://fun-dotto.com など。
func ArticleURL(baseURL, id string) string {
	return strings.TrimRight(baseURL, "/") + "/articles/" + url.PathEscape(id)
}

// ArticleBaseURLFromEnv は環境変数 ARTICLE_BASE_URL を返す。
func ArticleBaseURLFromEnv() (string, error) {
	v := os.Getenv("ARTICLE_BASE_URL")
	if v == "" {
		return "", errors.New("ARTICLE_BASE_URL is required")
	}
	return v, nil
}
