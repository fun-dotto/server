package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	loginURL           = "https://students.fun.ac.jp/Login"
	classChangeListURL = "https://students.fun.ac.jp/Pt/CSLecture"

	// loginPostDelay はログイン直後の遷移待ち。Python 版に合わせて 2 秒固定する。
	loginPostDelay = 2 * time.Second
)

// hiddenFieldNames は ASP.NET WebForms のログインフォームに含まれる隠しフィールド。
var hiddenFieldNames = []string{"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"}

// Client は学生ポータル（students.fun.ac.jp）へログインし、休講・補講・部屋変更の
// 一覧ページを取得する。
type Client struct {
	httpClient *http.Client
	year       int
}

// NewClient は対象年度 year（ログインフォームの TargetYearList に使う）を指定して Client を作る。
func NewClient(year int) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}
	return &Client{
		httpClient: &http.Client{Jar: jar, Timeout: 30 * time.Second},
		year:       year,
	}, nil
}

// FetchClassChangeHTML はログインしたのち /Pt/CSLecture の HTML を取得する。
func (c *Client) FetchClassChangeHTML(ctx context.Context, userID, password string) (string, error) {
	if err := c.login(ctx, userID, password); err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, classChangeListURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch class change list: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch class change list: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read class change list body: %w", err)
	}
	return string(body), nil
}

// login はログインページから隠しフィールドを取得し、認証情報とともに POST する。
func (c *Client) login(ctx context.Context, userID, password string) error {
	getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, loginURL, nil)
	if err != nil {
		return err
	}
	getResp, err := c.httpClient.Do(getReq)
	if err != nil {
		return fmt.Errorf("fetch login page: %w", err)
	}
	defer func() { _ = getResp.Body.Close() }()
	if getResp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch login page: unexpected status %d", getResp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(getResp.Body)
	if err != nil {
		return fmt.Errorf("parse login page: %w", err)
	}

	form := url.Values{}
	form.Set("__LASTFOCUS", "")
	form.Set("__EVENTTARGET", "")
	form.Set("__EVENTARGUMENT", "")
	form.Set("__SCROLLPOSITIONX", "0")
	form.Set("__SCROLLPOSITIONY", "0")
	form.Set("ctl00$MainContent$TargetYearList", strconv.Itoa(c.year))
	form.Set("ctl00$MainContent$TargetTermList", "11")
	form.Set("ctl00$MainContent$LoginId", userID)
	form.Set("ctl00$MainContent$LoginPassword", password)
	form.Set("ctl00$MainContent$LoginButton", "ログイン")
	for _, name := range hiddenFieldNames {
		if v, ok := doc.Find(fmt.Sprintf(`[name="%s"]`, name)).Attr("value"); ok {
			form.Set(name, v)
		}
	}

	postReq, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postResp, err := c.httpClient.Do(postReq)
	if err != nil {
		return fmt.Errorf("post login form: %w", err)
	}
	defer func() { _ = postResp.Body.Close() }()
	if _, err := io.Copy(io.Discard, postResp.Body); err != nil {
		return fmt.Errorf("drain login response: %w", err)
	}

	// ポータル側のセッション確立を待つ。Python 版と同じ固定 2 秒。
	select {
	case <-time.After(loginPostDelay):
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
