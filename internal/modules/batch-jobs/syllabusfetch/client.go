package syllabusfetch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

const (
	listURL   = "https://students.fun.ac.jp/lesson/SyllabusList"
	detailURL = "https://students.fun.ac.jp/lesson/Syllabus"

	yearFieldName   = "ctl00$MainContent$param_Syllabus_year_eq"
	searchFieldName = "ctl00$MainContent$SearchButton"
	gridViewTarget  = "ctl00$MainContent$GridView"
)

// Client は教務システムのシラバスページを取得する。
// 教務システムへの負荷を抑えるため、リクエストごとに Interval だけ間隔を空ける。
type Client struct {
	http     *http.Client
	interval time.Duration
	last     time.Time
}

func NewClient(interval time.Duration) *Client {
	// 一覧のページングは ASP.NET のセッションに依存しうるため Cookie を保持する
	jar, _ := cookiejar.New(nil)
	return &Client{
		http:     &http.Client{Jar: jar, Timeout: 60 * time.Second},
		interval: interval,
	}
}

// ListLessonIDs は指定年度の一覧を全ページ辿り、授業コードを返す。
// 一覧は ASP.NET WebForms のため、検索・ページングとも直前のレスポンスの hidden 項目を付けて POST する。
func (c *Client) ListLessonIDs(ctx context.Context, year int) ([]string, error) {
	initial, err := c.fetchList(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch search form: %w", err)
	}

	form := formFrom(initial.HiddenFields, year)
	form.Set(searchFieldName, "検索の実行")
	page, err := c.fetchList(ctx, form)
	if err != nil {
		return nil, fmt.Errorf("search year %d: %w", year, err)
	}

	var ids []string
	seen := map[string]bool{}
	for current := 1; ; current++ {
		for _, id := range page.LessonIDs {
			if !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
		next := current + 1
		if !page.HasPage(next) {
			return ids, nil
		}
		form := formFrom(page.HiddenFields, year)
		form.Set("__EVENTTARGET", gridViewTarget)
		form.Set("__EVENTARGUMENT", "Page$"+strconv.Itoa(next))
		if page, err = c.fetchList(ctx, form); err != nil {
			return nil, fmt.Errorf("fetch list page %d: %w", next, err)
		}
	}
}

// FetchDetail は詳細ページの HTML を返す。
func (c *Client) FetchDetail(ctx context.Context, lessonID string, year int) ([]byte, error) {
	q := url.Values{"lesson_id": {lessonID}, "year": {strconv.Itoa(year)}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailURL+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) fetchList(ctx context.Context, form url.Values) (ListPage, error) {
	method, body := http.MethodGet, io.Reader(nil)
	if form != nil {
		method, body = http.MethodPost, bytes.NewBufferString(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, listURL, body)
	if err != nil {
		return ListPage{}, err
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	b, err := c.do(req)
	if err != nil {
		return ListPage{}, err
	}
	return ParseListPage(bytes.NewReader(b))
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	if wait := c.interval - time.Since(c.last); wait > 0 {
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(wait):
		}
	}
	defer func() { c.last = time.Now() }()

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s: unexpected status %d", req.Method, req.URL, res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

func formFrom(hidden map[string]string, year int) url.Values {
	form := url.Values{}
	for k, v := range hidden {
		form.Set(k, v)
	}
	form.Set(yearFieldName, strconv.Itoa(year))
	return form
}
