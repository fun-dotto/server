package syllabusfetch

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// listHTML は授業コードとページャーのリンクだけを持つ一覧ページを返す。
func listHTML(ids []string, pages []int) string {
	var b strings.Builder
	b.WriteString(`<html><body><form><input type="hidden" name="__VIEWSTATE" value="vs" /><div>`)
	for _, p := range pages {
		fmt.Fprintf(&b, `<a href="javascript:__doPostBack('ctl00$MainContent$GridView','Page$%d')">%d</a>`, p, p)
	}
	b.WriteString(`</div><table>`)
	for _, id := range ids {
		fmt.Fprintf(&b, `<tr><td><a href="Syllabus.aspx?lesson_id=%s&amp;year=2026">詳細</a></td></tr>`, id)
	}
	b.WriteString(`</table></form></body></html>`)
	return b.String()
}

// newListServer は検索で 1 ページ目、Page$N のポストバックで pages[N] を返す一覧のテストサーバーを立てる。
func newListServer(t *testing.T, pages map[string]string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/SyllabusList" {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			fmt.Fprint(w, listHTML(nil, nil))
			return
		}
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "2026", r.PostForm.Get(yearFieldName))
		assert.Equal(t, "vs", r.PostForm.Get("__VIEWSTATE"))
		key := "search"
		if arg := r.PostForm.Get("__EVENTARGUMENT"); arg != "" {
			key = arg
		}
		fmt.Fprint(w, pages[key])
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestClientListLessonIDs(t *testing.T) {
	t.Run("ページャーを辿って全ページの授業コードを返す", func(t *testing.T) {
		srv := newListServer(t, map[string]string{
			"search": listHTML([]string{"100101", "100201"}, []int{2, 3}),
			"Page$2": listHTML([]string{"100301"}, []int{1, 3}),
			"Page$3": listHTML([]string{"100401"}, []int{1, 2}),
		})
		got, err := NewClient(srv.URL, 0).ListLessonIDs(context.Background(), 2026)
		require.NoError(t, err)
		assert.Equal(t, []string{"100101", "100201", "100301", "100401"}, got)
	})

	t.Run("要求と違うページが返り新しい授業コードがなければエラー", func(t *testing.T) {
		first := listHTML([]string{"100101"}, []int{2})
		srv := newListServer(t, map[string]string{
			"search": first,
			// ViewState の期限切れなどで 1 ページ目が返った場合
			"Page$2": first,
		})
		_, err := NewClient(srv.URL, 0).ListLessonIDs(context.Background(), 2026)
		assert.Error(t, err)
	})
}
