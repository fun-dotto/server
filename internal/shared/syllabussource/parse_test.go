package syllabussource

import "testing"

const detailHTML = `<html><body><form>
<div class="card m20"><header class="card-header"><p class="card-header-title">基本情報</p></header>
<div class="card-content"><table>
<tr><th>授業名</th><td colspan="3"><span>認知科学1～4</span></td></tr>
<tr><th>担当教員名</th><td><span>◎花田　光彦<br>伊藤　精英</span></td></tr>
<tr><th>実務家教員区分</th><td><div>
   -
</div></td><th>複数人担当形式</th><td>オムニバス</td></tr>
<tr><th>対象コース・領域　※</th><td><div>全コース
<br />学部共通選択</div></td></tr>
</table></div></div>
<div class="card m20 is-hidden"><header><p class="card-header-title">ディプロマポリシーの各領域と配点</p></header>
<div class="card-content"><table><tr><th>隠し</th><td>x</td></tr></table></div></div>
<div class="card m20"><header><p class="card-header-title">
  授業内容とスケジュール
</p></header><div class="card-content"><div>1. A
<br />2. B</div></div></div>
</form></body></html>`

func TestParseDetailPage(t *testing.T) {
	d, err := parseDetailPage(detailHTML)
	if err != nil {
		t.Fatal(err)
	}
	want := Detail{
		"授業名":         "認知科学1～4",
		"担当教員名":       "◎花田　光彦\n伊藤　精英",
		"実務家教員区分":     "-",
		"複数人担当形式":     "オムニバス",
		"対象コース・領域":    "全コース\n学部共通選択",
		"授業内容とスケジュール": "1. A\n2. B",
	}
	if len(d) != len(want) {
		t.Errorf("got %d labels: %v", len(d), d)
	}
	for k, v := range want {
		if d[k] != v {
			t.Errorf("%s = %q, want %q", k, d[k], v)
		}
	}
	m, unknown := ToModel(2026, ListEntry{LessonID: "100101", DisplayOrder: "0"}, d)
	if len(unknown) != 0 || m.TargetCourses != want["対象コース・領域"] || m.LessonID != "100101" || m.Year != 2026 {
		t.Errorf("ToModel = %+v, unknown %v", m, unknown)
	}
}

const listHTML = `<h3>検索結果(21)</h3><table id="MainContent_GridView">
<tr class="Pager"><td><table><tr><td>1</td></tr></table></td></tr>
<tr><th>履修年度</th></tr>
<tr><td><a href="#">詳細</a></td><td data-col-responsive-title="履修年度">2026</td><td data-col-responsive-title="授業コード">100101</td><td data-col-responsive-title="授業科目名">認知科学1～4</td><td data-col-responsive-title="Course">Cognitive Science</td><td data-col-responsive-title="担当教員名">花田　光彦</td><td data-col-responsive-title="授業形態">講義</td><td data-col-responsive-title="表示順">0</td></tr>
</table>`

func TestParseListPage(t *testing.T) {
	r, err := parseListPage(listHTML)
	if err != nil {
		t.Fatal(err)
	}
	if r.total != 21 || len(r.entries) != 1 {
		t.Fatalf("got %+v", r)
	}
	want := ListEntry{Year: "2026", LessonID: "100101", LessonName: "認知科学1～4", EnName: "Cognitive Science", StaffName: "花田　光彦", TeachingForm: "講義", DisplayOrder: "0"}
	if r.entries[0] != want {
		t.Errorf("got %+v", r.entries[0])
	}
}
