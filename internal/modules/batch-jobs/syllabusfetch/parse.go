package syllabusfetch

import (
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/fun-dotto/server/internal/shared/model"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ListPage は一覧ページ 1 ページ分の解析結果。
type ListPage struct {
	// LessonIDs はこのページに載っている授業コード（表示順）。
	LessonIDs []string
	// HiddenFields は次のポストバックにそのまま送り返す ASP.NET の hidden 項目（__VIEWSTATE など）。
	HiddenFields map[string]string
	// PageNumbers はページャーに表示されているページ番号。表示は前後数ページに限られるため全ページ数ではない。
	PageNumbers []int
}

// HasPage は指定したページへのリンクがページャーにあるかを返す。
func (p ListPage) HasPage(n int) bool {
	return slices.Contains(p.PageNumbers, n)
}

var (
	detailLinkPattern = regexp.MustCompile(`lesson_id=([^&]+)`)
	pagerPattern      = regexp.MustCompile(`Page\$(\d+)`)
	sourceNewline     = strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ")
)

// ParseListPage は一覧ページ（SyllabusList）の HTML を解析する。
func ParseListPage(r io.Reader) (ListPage, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return ListPage{}, err
	}
	page := ListPage{HiddenFields: map[string]string{}}
	seenPage := map[int]bool{}
	for n := range doc.Descendants() {
		if n.Type != html.ElementNode {
			continue
		}
		switch n.DataAtom {
		case atom.Input:
			if attr(n, "type") == "hidden" && attr(n, "name") != "" {
				page.HiddenFields[attr(n, "name")] = attr(n, "value")
			}
		case atom.A:
			href := attr(n, "href")
			if strings.Contains(href, "Syllabus.aspx?") || strings.Contains(href, "Syllabus?") {
				if m := detailLinkPattern.FindStringSubmatch(href); m != nil {
					page.LessonIDs = append(page.LessonIDs, m[1])
				}
				continue
			}
			for _, m := range pagerPattern.FindAllStringSubmatch(href, -1) {
				num, err := strconv.Atoi(m[1])
				if err == nil && !seenPage[num] {
					seenPage[num] = true
					page.PageNumbers = append(page.PageNumbers, num)
				}
			}
		}
	}
	slices.Sort(page.PageNumbers)
	return page, nil
}

// ParseDetailPage は詳細ページ（Syllabus）の HTML から見出しごとの値を取り出し、Syllabus に詰める。
// LessonID・Year・ContentHash は呼び出し側で設定する。
// 対応する列のない見出しは unknownLabels として返す（取得元に項目が増えたことに気づくため）。
func ParseDetailPage(r io.Reader) (s model.Syllabus, unknownLabels []string, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return model.Syllabus{}, nil, err
	}
	for label, value := range detailItems(doc) {
		field := syllabusField(&s, label)
		if field == nil {
			unknownLabels = append(unknownLabels, label)
			continue
		}
		*field = value
	}
	slices.Sort(unknownLabels)
	return s, unknownLabels, nil
}

// detailItems は表の <th> と直後の <td>、および表を持たないカード（授業内容とスケジュール）の見出しと本文を返す。
// 非表示（is-hidden）の要素は科目によって使われない枠のため読まない。
func detailItems(doc *html.Node) map[string]string {
	items := map[string]string{}
	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && hasClass(n, "is-hidden") {
			return
		}
		if n.Type == html.ElementNode {
			switch {
			case n.DataAtom == atom.Tr:
				collectRow(n, items)
				return
			case hasClass(n, "card") && findFirst(n, func(c *html.Node) bool { return c.DataAtom == atom.Table }) == nil:
				title := findFirst(n, func(c *html.Node) bool { return hasClass(c, "card-header-title") })
				content := findFirst(n, func(c *html.Node) bool { return hasClass(c, "card-content") })
				if title != nil && content != nil {
					items[normalizeLabel(textOf(title))] = textOf(content)
				}
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return items
}

// collectRow は 1 行に「見出し・値」が複数組並ぶ行（実務家教員区分と複数人担当形式など）も扱う。
func collectRow(tr *html.Node, items map[string]string) {
	label := ""
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode || hasClass(c, "is-hidden") {
			continue
		}
		switch c.DataAtom {
		case atom.Th:
			label = normalizeLabel(textOf(c))
		case atom.Td:
			if label != "" {
				items[label] = textOf(c)
				label = ""
			}
		}
	}
}

// syllabusField は見出しに対応する Syllabus の列を返す。
// 「予定開講期」は開講期が確定していない科目で「開講期」の代わりに表示される。
func syllabusField(s *model.Syllabus, label string) *string {
	switch label {
	case "授業名":
		return &s.Name
	case "Course":
		return &s.EnName
	case "配当年次":
		return &s.Grades
	case "開講期", "予定開講期":
		return &s.Term
	case "単位数":
		return &s.Credit
	case "担当教員名":
		return &s.FacultyNames
	case "実務家教員区分":
		return &s.PracticalHomeFacultyCategory
	case "複数人担当形式":
		return &s.MultiplePersonTeachingForm
	case "授業形態":
		return &s.TeachingForm
	case "授業の概要":
		return &s.Summary
	case "授業の到達目標":
		return &s.LearningOutcomes
	case "提出課題等":
		return &s.Assignments
	case "成績の評価方法・基準":
		return &s.EvaluationMethod
	case "テキスト":
		return &s.Textbooks
	case "参考書":
		return &s.ReferenceBooks
	case "履修条件":
		return &s.Prerequisites
	case "事前学習":
		return &s.PreLearning
	case "事後学習":
		return &s.PostLearning
	case "履修上の留意点":
		return &s.Notes
	case "キーワード/keywords":
		return &s.Keywords
	case "対象コース・領域":
		return &s.TargetCourses
	case "科目群・科目区分":
		return &s.Classifications
	case "DSOP対象科目":
		return &s.DsopSubject
	case "教授言語":
		return &s.TeachingLanguage
	case "授業・試験の形式":
		return &s.TeachingAndExamForm
	case "授業内容とスケジュール":
		return &s.ContentsAndSchedule
	case "(E)Course outline":
		return &s.EnSummary
	case "(E)Course Objectives":
		return &s.EnLearningOutcomes
	case "(E)Course Schedule":
		return &s.EnContentsAndSchedule
	case "(E)Assessment":
		return &s.EnEvaluationMethod
	case "(E)Language of Instruction":
		return &s.EnTeachingLanguage
	case "(E)Requirements for registration":
		return &s.EnPrerequisites
	case "(E)Prior/Post Assignment":
		return &s.EnPrePostLearning
	case "(E)Textbooks":
		return &s.EnTextbooks
	case "(E)Note":
		return &s.EnNotes
	case "(E)Type of class and exam":
		return &s.EnTeachingAndExamForm
	case "ディプロマポリシー":
		return &s.DiplomaPolicies
	case "ディプロマポリシーとの対応":
		return &s.DiplomaPolicyCorrespondence
	}
	return nil
}

// normalizeLabel は見出しの装飾（全角空白と「※」）を除く。例: 「対象コース・領域　※」
func normalizeLabel(s string) string {
	s = strings.ReplaceAll(s, "※", "")
	return strings.TrimSpace(strings.ReplaceAll(s, "　", " "))
}

// textOf は要素の表示テキストを返す。HTML のインデント由来の行頭・行末の空白は除く。
// 全角空白は氏名の区切りなど内容の一部のため残す。
func textOf(n *html.Node) string {
	var b strings.Builder
	for d := range n.Descendants() {
		switch {
		case d.Type == html.TextNode:
			// ソース上の改行はブラウザ同様に空白として扱い、改行は <br> だけで表す
			b.WriteString(sourceNewline.Replace(d.Data))
		case d.Type == html.ElementNode && d.DataAtom == atom.Br:
			b.WriteString("\n")
		}
	}
	lines := strings.Split(b.String(), "\n")
	for i, l := range lines {
		lines[i] = strings.Trim(l, " \t\r")
	}
	return strings.Trim(strings.Join(lines, "\n"), "\n")
}

func findFirst(n *html.Node, match func(*html.Node) bool) *html.Node {
	for d := range n.Descendants() {
		if d.Type == html.ElementNode && match(d) {
			return d
		}
	}
	return nil
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
	return slices.Contains(strings.Fields(attr(n, "class")), class)
}
