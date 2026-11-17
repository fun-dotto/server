package syllabusfetch

import (
	"strings"
	"testing"

	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 教務システムの一覧ページ（2026-11 時点）の構造を縮めたもの。
const listPageHTML = `<html><body>
<form method="post" action="./SyllabusList" id="Form1">
<input type="hidden" name="__VIEWSTATE" id="__VIEWSTATE" value="vs+/=" />
<input type="hidden" name="__EVENTVALIDATION" id="__EVENTVALIDATION" value="ev" />
<input name="ctl00$MainContent$param_lesson_name_like" type="text" />
<div class="level"><a class="button is-link" href="javascript:__doPostBack('ctl00$MainContent$GridView','Page$1')">前へ</a>
<a class="button" href="javascript:__doPostBack('ctl00$MainContent$GridView','Page$3')">3</a>
<a class="button" href="javascript:__doPostBack('ctl00$MainContent$GridView','Page$18')">18</a></div>
<table>
<tr class="Pager is-hidden"><td><table><tr><td><span>2</span></td>
<td><a href="javascript:__doPostBack(&#39;ctl00$MainContent$GridView&#39;,&#39;Page$Last&#39;)">&gt;&gt;</a></td></tr></table></td></tr>
<tr><th scope="col">&nbsp;</th><th scope="col">授業コード</th></tr>
<tr><td><a class="button is-link" href="Syllabus.aspx?lesson_id=100101&amp;year=2026&amp;from=SyllabusList.aspx">詳細</a></td><td>100101</td></tr>
<tr><td><a class="button is-link" href="Syllabus.aspx?lesson_id=100201&amp;year=2026&amp;from=SyllabusList.aspx">詳細</a></td><td>100201</td></tr>
</table>
</form></body></html>`

func TestParseListPage(t *testing.T) {
	got, err := ParseListPage(strings.NewReader(listPageHTML))
	require.NoError(t, err)
	want := ListPage{
		LessonIDs:    []string{"100101", "100201"},
		HiddenFields: map[string]string{"__VIEWSTATE": "vs+/=", "__EVENTVALIDATION": "ev"},
		PageNumbers:  []int{1, 3, 18},
	}
	assert.Equal(t, want, got)
}

func TestListPageHasPage(t *testing.T) {
	page := ListPage{PageNumbers: []int{1, 3, 18}}
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{name: "ページャーにある", n: 3, want: true},
		{name: "ページャーにない（最終ページの次）", n: 19, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, page.HasPage(tt.n))
		})
	}
}

// 教務システムの詳細ページ（2026-11 時点）の構造を縮めたもの。実際のページは <br /> の後にソース上の改行（CRLF）が続く。
const detailPageHTML = `<html><body>
<div class="card m20">
  <header class="card-header js-card-toggle"><p class="card-header-title">基本情報</p></header>
  <div class="card-content"><div class="custom-table-width"><table class="table">
    <tr>
      <th>授業名</th>
      <td colspan="3">
        <span id="MainContent_lesson_name" class="aspNetDisabled">認知科学1～4</span>
      </td>
    </tr>
    <tr id="MainContent_display_linked_timetable_terms">
	<th>予定開講期</th>
	<td colspan="3"><span>前期</span></td>
    </tr>
    <tr>
      <th>担当教員名</th>
      <td colspan="3"><span>◎花田　光彦<br>伊藤　精英</span></td>
    </tr>
    <tr>
      <th>実務家教員区分</th>
      <td><div class="field"><div class="control">
          -
      </div></div></td>
      <th>複数人担当形式</th>
      <td><div class="field"><div class="control">
          オムニバス
      </div></div></td>
    </tr>
    <tr>
      <th colspan="1">事前学習
      </th>
      <td colspan="3"><div id='auto_self_study' class=''>1行目<br /><br />3行目</div></td>
    </tr>
    <tr>
      <th>対象コース・領域　※</th>
      <td colspan="3"><div id='auto_option2' class=''>全コース<br />学部共通選択</div></td>
    </tr>
    <tr>
      <th>キーワード/keywords</th>
      <td colspan="3"><div id='auto_option1' class=''>思考 / thinking<br /><br /></div></td>
    </tr>
    <tr>
      <th>(E)Course outline</th>
      <td colspan="3"><div id='auto_option6' class=''>In this lecture, &amp; more.</div></td>
    </tr>
    <tr class="is-hidden">
      <th>教授言語</th>
      <td colspan="3">日本語</td>
    </tr>
    <tr>
      <th>新しい項目</th>
      <td colspan="3">x</td>
    </tr>
  </table></div></div>
</div>
<div class="card m20 is-hidden">
  <header class="card-header"><p class="card-header-title">ディプロマポリシーの各領域と配点</p></header>
  <div class="card-content"><table><tr class="is-hidden"><td><p class="has-text-weight-bold"></p></td></tr></table></div>
</div>
<div class="card m20">
  <header class="card-header js-card-toggle">
    <p class="card-header-title">
        授業内容とスケジュール
    </p>
  </header>
  <div class="card-content"><div class="field"><div class="control">
    <div id='auto_syllabus_details' class=''>1. 導入<br />` + "\r\n" + `2. まとめ</div>
  </div></div></div>
</div>
</body></html>`

func TestParseDetailPage(t *testing.T) {
	got, unknown, err := ParseDetailPage(strings.NewReader(detailPageHTML))
	require.NoError(t, err)
	want := model.Syllabus{
		Name:                         "認知科学1～4",
		Term:                         "前期",
		FacultyNames:                 "◎花田　光彦\n伊藤　精英",
		PracticalHomeFacultyCategory: "-",
		MultiplePersonTeachingForm:   "オムニバス",
		PreLearning:                  "1行目\n\n3行目",
		TargetCourses:                "全コース\n学部共通選択",
		Keywords:                     "思考 / thinking",
		EnSummary:                    "In this lecture, & more.",
		ContentsAndSchedule:          "1. 導入\n2. まとめ",
	}
	assert.Equal(t, want, got)
	assert.Equal(t, []string{"新しい項目"}, unknown)
}
