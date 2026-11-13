package syllabussource

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ListEntry は一覧ページの 1 行。値は表示されている文字列のまま。
type ListEntry struct {
	Year         string // 履修年度
	LessonID     string // 授業コード
	LessonName   string // 授業科目名
	EnName       string // Course
	StaffName    string // 担当教員名（主担当のみ）
	TeachingForm string // 授業形態
	DisplayOrder string // 表示順
}

// Detail は詳細ページの見出し（正規化済み）と値の対応。
type Detail map[string]string

type listResult struct {
	entries []ListEntry
	total   int
}

var totalPattern = regexp.MustCompile(`検索結果\((\d+)\)`)

func parseListPage(page string) (listResult, error) {
	m := totalPattern.FindStringSubmatch(page)
	if m == nil {
		return listResult{}, fmt.Errorf("result count not found")
	}
	total, _ := strconv.Atoi(m[1])

	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return listResult{}, err
	}
	table := find(doc, func(n *html.Node) bool { return n.DataAtom == atom.Table && attr(n, "id") == "MainContent_GridView" })
	if table == nil {
		if total == 0 {
			return listResult{total: 0}, nil
		}
		return listResult{}, fmt.Errorf("result table not found")
	}
	var entries []ListEntry
	for _, tr := range findAll(table, func(n *html.Node) bool { return n.DataAtom == atom.Tr }) {
		cells := map[string]string{}
		for c := tr.FirstChild; c != nil; c = c.NextSibling {
			if c.DataAtom == atom.Td {
				if title := attr(c, "data-col-responsive-title"); title != "" {
					cells[title] = text(c)
				}
			}
		}
		if cells["授業コード"] == "" {
			continue
		}
		entries = append(entries, ListEntry{
			Year:         cells["履修年度"],
			LessonID:     cells["授業コード"],
			LessonName:   cells["授業科目名"],
			EnName:       cells["Course"],
			StaffName:    cells["担当教員名"],
			TeachingForm: cells["授業形態"],
			DisplayOrder: cells["表示順"],
		})
	}
	return listResult{entries: entries, total: total}, nil
}

// hiddenFields はポストバックに必要な hidden input（__VIEWSTATE など）を集める。
func hiddenFields(page string) (url.Values, error) {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return nil, err
	}
	form := find(doc, func(n *html.Node) bool { return n.DataAtom == atom.Form })
	if form == nil {
		return nil, errNoForm
	}
	v := url.Values{}
	for _, in := range findAll(form, func(n *html.Node) bool {
		return n.DataAtom == atom.Input && strings.EqualFold(attr(n, "type"), "hidden")
	}) {
		if name := attr(in, "name"); name != "" {
			v.Set(name, attr(in, "value"))
		}
	}
	if v.Get("__VIEWSTATE") == "" {
		return nil, errNoForm
	}
	return v, nil
}

// parseDetailPage は表示されているカード（is-hidden を除く）から見出しと値を取り出す。
// 表の見出しは <th>、表を持たないカード（授業内容とスケジュールなど）はカードの見出しを使う。
func parseDetailPage(page string) (Detail, error) {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return nil, err
	}
	d := Detail{}
	cards := findAll(doc, func(n *html.Node) bool { return n.DataAtom == atom.Div && hasClass(n, "card") })
	if len(cards) == 0 {
		return nil, fmt.Errorf("syllabus not found")
	}
	for _, card := range cards {
		if hidden(card) {
			continue
		}
		content := find(card, func(n *html.Node) bool { return hasClass(n, "card-content") })
		if content == nil {
			continue
		}
		ths := findAll(content, func(n *html.Node) bool { return n.DataAtom == atom.Th && !hidden(n) })
		if len(ths) == 0 {
			title := find(card, func(n *html.Node) bool { return hasClass(n, "card-header-title") })
			if title != nil {
				d.set(NormalizeLabel(text(title)), text(content))
			}
			continue
		}
		for _, th := range ths {
			td := th.NextSibling
			for td != nil && td.DataAtom != atom.Td {
				td = td.NextSibling
			}
			if td == nil {
				continue
			}
			d.set(NormalizeLabel(text(th)), text(td))
		}
	}
	if len(d) == 0 {
		return nil, fmt.Errorf("syllabus not found")
	}
	return d, nil
}

func (d Detail) set(label, value string) {
	if label == "" {
		return
	}
	if old, ok := d[label]; ok && old != "" && value == "" {
		return
	}
	d[label] = value
}

// NormalizeLabel は見出しの前後の空白と末尾の「※」を取り除く。
func NormalizeLabel(s string) string {
	s = strings.TrimSpace(s)
	for {
		t := strings.TrimSpace(strings.TrimSuffix(s, "※"))
		if t == s {
			return t
		}
		s = t
	}
}

// text はブラウザでの表示に近い文字列を返す。
// HTML 上の空白・改行（ASCII）は 1 つの空白にまとめ、<br> を改行にし、各行の前後の空白を取り除く。
func text(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			b.WriteString(collapseSpace(n.Data))
			return
		case html.ElementNode:
			switch n.DataAtom {
			case atom.Br:
				b.WriteString("\n")
				return
			case atom.Script, atom.Style:
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	lines := strings.Split(b.String(), "\n")
	for i, l := range lines {
		lines[i] = strings.Trim(l, " ")
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n")
}

func collapseSpace(s string) string {
	var b strings.Builder
	space := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' {
			if !space {
				b.WriteByte(' ')
			}
			space = true
			continue
		}
		space = false
		b.WriteRune(r)
	}
	return b.String()
}

func hidden(n *html.Node) bool {
	for p := n; p != nil; p = p.Parent {
		if hasClass(p, "is-hidden") {
			return true
		}
	}
	return false
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func hasClass(n *html.Node, class string) bool {
	if n.Type != html.ElementNode {
		return false
	}
	for _, c := range strings.Fields(attr(n, "class")) {
		if c == class {
			return true
		}
	}
	return false
}

func find(n *html.Node, match func(*html.Node) bool) *html.Node {
	if n.Type == html.ElementNode && match(n) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if f := find(c, match); f != nil {
			return f
		}
	}
	return nil
}

func findAll(n *html.Node, match func(*html.Node) bool) []*html.Node {
	var out []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && match(n) {
			out = append(out, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return out
}
