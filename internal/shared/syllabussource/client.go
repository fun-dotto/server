// Package syllabussource は教務システムのシラバス（一覧・詳細ページ）を取得する。
// ログインは不要で、一覧は ASP.NET WebForms のポストバックで検索・ページングする。
package syllabussource

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://students.fun.ac.jp/lesson/"
	listPath       = "SyllabusList"
	detailPath     = "Syllabus"

	yearField      = "ctl00$MainContent$param_Syllabus_year_eq"
	searchButton   = "ctl00$MainContent$SearchButton"
	gridViewTarget = "ctl00$MainContent$GridView"
	// listPageSize は一覧の 1 ページあたりの件数。
	listPageSize = 20
	// maxListPages は一覧のページ数の上限（無限ループ防止）。
	maxListPages = 200
)

// Client は教務システムのシラバスページを取得する。
// リクエストの間隔を Interval 以上空ける。
type Client struct {
	BaseURL   string
	Interval  time.Duration
	UserAgent string

	http *http.Client
	last time.Time
}

func NewClient(interval time.Duration) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		BaseURL:   DefaultBaseURL,
		Interval:  interval,
		UserAgent: "dotto-syllabus-fetcher",
		http:      &http.Client{Timeout: 60 * time.Second, Jar: jar},
	}
}

// ListAll は year の一覧を全ページ取得する。
func (c *Client) ListAll(ctx context.Context, year int) ([]ListEntry, error) {
	page, err := c.get(ctx, c.BaseURL+listPath)
	if err != nil {
		return nil, fmt.Errorf("get list page: %w", err)
	}
	form, err := hiddenFields(page)
	if err != nil {
		return nil, err
	}
	form.Set(yearField, strconv.Itoa(year))
	form.Set(searchButton, "検索の実行")
	page, err = c.post(ctx, c.BaseURL+listPath, form)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	result, err := parseListPage(page)
	if err != nil {
		return nil, fmt.Errorf("parse list page 1: %w", err)
	}
	entries := result.entries
	total := result.total
	for p := 2; len(entries) < total; p++ {
		if p > maxListPages || p > (total+listPageSize-1)/listPageSize {
			return nil, fmt.Errorf("list: got %d entries, want %d", len(entries), total)
		}
		form, err := hiddenFields(page)
		if err != nil {
			return nil, err
		}
		form.Set("__EVENTTARGET", gridViewTarget)
		form.Set("__EVENTARGUMENT", "Page$"+strconv.Itoa(p))
		page, err = c.post(ctx, c.BaseURL+listPath, form)
		if err != nil {
			return nil, fmt.Errorf("list page %d: %w", p, err)
		}
		result, err := parseListPage(page)
		if err != nil {
			return nil, fmt.Errorf("parse list page %d: %w", p, err)
		}
		if len(result.entries) == 0 {
			return nil, fmt.Errorf("list page %d is empty", p)
		}
		entries = append(entries, result.entries...)
	}
	for _, e := range entries {
		if e.Year != strconv.Itoa(year) {
			return nil, fmt.Errorf("list: unexpected year %q for lesson %s", e.Year, e.LessonID)
		}
	}
	return entries, nil
}

// Detail は授業の詳細ページを取得する。
func (c *Client) Detail(ctx context.Context, lessonID string, year int) (Detail, error) {
	q := url.Values{"lesson_id": {lessonID}, "year": {strconv.Itoa(year)}}
	page, err := c.get(ctx, c.BaseURL+detailPath+"?"+q.Encode())
	if err != nil {
		return nil, err
	}
	return parseDetailPage(page)
}

func (c *Client) get(ctx context.Context, u string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	return c.do(req)
}

func (c *Client) post(ctx context.Context, u string, form url.Values) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req)
}

func (c *Client) do(req *http.Request) (string, error) {
	if err := c.wait(req.Context()); err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	res, err := c.http.Do(req)
	c.last = time.Now()
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s %s: status %d", req.Method, req.URL, res.StatusCode)
	}
	return string(body), nil
}

func (c *Client) wait(ctx context.Context) error {
	if c.last.IsZero() {
		return nil
	}
	d := c.Interval - time.Since(c.last)
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

var errNoForm = errors.New("hidden fields not found")
