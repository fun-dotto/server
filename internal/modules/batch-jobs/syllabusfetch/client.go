package syllabusfetch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

const (
	// DefaultBaseURL は教務システムのシラバスページのベース URL。
	DefaultBaseURL = "https://students.fun.ac.jp/lesson"

	yearFieldName   = "ctl00$MainContent$param_Syllabus_year_eq"
	searchFieldName = "ctl00$MainContent$SearchButton"
	gridViewTarget  = "ctl00$MainContent$GridView"
)

// Client は教務システムのシラバスページを取得する。
// 教務システムへの負荷を抑えるため、リクエストごとに Interval だけ間隔を空ける。
type Client struct {
	baseURL  string
	http     *http.Client
	interval time.Duration
	last     time.Time
}

// NewClient は baseURL（末尾の / なし）配下の SyllabusList・Syllabus を取得するクライアントを返す。
func NewClient(baseURL string, interval time.Duration) *Client {
	// 一覧のページングは ASP.NET のセッションに依存しうるため Cookie を保持する
	jar, _ := cookiejar.New(nil)
	return &Client{
		baseURL:  baseURL,
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
		added := 0
		for _, id := range page.LessonIDs {
			if !seen[id] {
				seen[id] = true
				ids = append(ids, id)
				added++
			}
		}
		// ViewState の期限切れなどで要求と違うページが返ると、既出の授業コードだけのページや 0 件のページになる。
		// そのまま終えると一部の科目を取りこぼしたまま成功扱いになるため、エラーにする。
		if added == 0 {
			return nil, fmt.Errorf("list page %d has no new lesson ids", current)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/Syllabus?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) fetchList(ctx context.Context, form url.Values) (listPage, error) {
	method, body := http.MethodGet, io.Reader(nil)
	if form != nil {
		method, body = http.MethodPost, bytes.NewBufferString(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+"/SyllabusList", body)
	if err != nil {
		return listPage{}, err
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	b, err := c.do(req)
	if err != nil {
		return listPage{}, err
	}
	return parseListPage(bytes.NewReader(b))
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
		return nil, &statusError{Method: req.Method, URL: req.URL.String(), Code: res.StatusCode}
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

// statusError は 200 以外のステータスが返ったことを表す。
type statusError struct {
	Method string
	URL    string
	Code   int
}

func (e *statusError) Error() string {
	return fmt.Sprintf("%s %s: unexpected status %d", e.Method, e.URL, e.Code)
}

// isRetryable は一時的な失敗（5xx・タイムアウト）かを返す。
func isRetryable(err error) bool {
	var se *statusError
	if errors.As(err, &se) {
		return se.Code >= 500
	}
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}
