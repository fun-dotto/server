// Package subjectbuild はシラバス（取得元の写し）から科目の構造化データを生成する。
// 科目生成ジョブと admin-api（シラバスの付け替え時）で使う。
package subjectbuild

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/fun-dotto/server/internal/shared/model"
)

const (
	courseInformationSystem = "InformationSystem"
	courseInformationDesign = "InformationDesign"
	courseAdvancedICT       = "AdvancedICT"
	courseComplexSystem     = "ComplexSystem"
	courseIntelligentSystem = "IntelligentSystem"
)

var allCourses = []string{courseInformationSystem, courseInformationDesign, courseAdvancedICT, courseComplexSystem, courseIntelligentSystem}

// Result は生成した科目と、生成時の警告。
type Result struct {
	Subject model.Subject
	// FacultyNames は担当教員名（◎ を除いた表示名）と主担当かどうか。Faculties の解決に使う。
	FacultyNames []FacultyName
	Warnings     []string
}

type FacultyName struct {
	Name      string
	IsPrimary bool
}

// Build は s から科目を生成する。Faculties は ResolveFaculties で設定する。
// 開講期が解釈できないなど、科目として必須の値が欠ける場合はエラーを返す。
func Build(s model.Syllabus) (Result, error) {
	var r Result
	semester := Semester(s.Term, s.Name)
	if semester == "" {
		return r, fmt.Errorf("開講期を解釈できません: %q", s.Term)
	}
	credit, ok := Credit(s.Credit)
	if !ok {
		return r, fmt.Errorf("単位数を解釈できません: %q", s.Credit)
	}
	classification, category := Classification(s.Classifications)
	attrs, warn := EligibleAttributes(s.Name, s.Grades, isGraduate(s))
	if warn != "" {
		r.Warnings = append(r.Warnings, warn)
	}
	r.Subject = model.Subject{
		Name:                    s.Name,
		Year:                    s.Year,
		Semester:                semester,
		Credit:                  credit,
		Classification:          classification,
		CulturalSubjectCategory: category,
		SyllabusID:              s.ID,
		EligibleAttributes:      attrs,
		Requirements:            Requirements(s.TargetCourses, s.Classifications, isGraduate(s)),
	}
	r.FacultyNames = FacultyNames(s.FacultyNames)
	return r, nil
}

// ResolveFaculties は担当教員名を faculties と名前で突合し、Subject.Faculties を設定する。
// 見つからない教員は含めず、名前を返す。
func (r *Result) ResolveFaculties(faculties []model.Faculty) (missing []string) {
	byName := make(map[string]model.Faculty, len(faculties))
	for _, f := range faculties {
		byName[normalizeName(f.Name)] = f
	}
	r.Subject.Faculties = nil
	seen := map[string]bool{}
	for _, n := range r.FacultyNames {
		f, ok := byName[normalizeName(n.Name)]
		if !ok {
			missing = append(missing, n.Name)
			continue
		}
		if seen[f.ID.String()] {
			continue
		}
		seen[f.ID.String()] = true
		r.Subject.Faculties = append(r.Subject.Faculties, model.SubjectFaculty{FacultyID: f.ID, IsPrimary: n.IsPrimary})
	}
	return missing
}

// Semester は開講期（前期/後期/1Q〜4Q/通年）を enum に変換する。授業名に夏期・冬期集中とあれば集中講義とする。
func Semester(term, name string) string {
	switch {
	case strings.Contains(name, "夏期集中"):
		return "SummerIntensive"
	case strings.Contains(name, "冬期集中"):
		return "WinterIntensive"
	}
	switch strings.TrimSpace(toHalfWidth(term)) {
	case "前期":
		return "H1"
	case "後期":
		return "H2"
	case "1Q":
		return "Q1"
	case "2Q":
		return "Q2"
	case "3Q":
		return "Q3"
	case "4Q":
		return "Q4"
	case "通年":
		return "AllYear"
	}
	return ""
}

var numberPattern = regexp.MustCompile(`\d+`)

// Credit は「2単位」を 2 に変換する。「0単位～2単位」のような範囲は最大値を使う。
func Credit(s string) (int, bool) {
	nums := numberPattern.FindAllString(toHalfWidth(s), -1)
	if len(nums) == 0 {
		return 0, false
	}
	max := 0
	for _, n := range nums {
		v, _ := strconv.Atoi(n)
		max = maxInt(max, v)
	}
	return max, true
}

// Classification は科目群・科目区分から科目カテゴリと教養科目カテゴリを決める。
func Classification(classifications string) (classification, category string) {
	c := classifications
	switch {
	case strings.Contains(c, "研究指導"):
		return "ResearchInstruction", ""
	case strings.Contains(c, "コミュニケーション科目群"):
		return "Cultural", "Communication"
	case strings.Contains(c, "教養"):
		switch {
		case strings.Contains(c, "社会"):
			category = "Society"
		case strings.Contains(c, "人間"):
			category = "Human"
		case strings.Contains(c, "科学"):
			category = "Science"
		case strings.Contains(c, "健康"):
			category = "Health"
		case strings.Contains(c, "コミュニケーション"):
			category = "Communication"
		}
		return "Cultural", category
	}
	return "Specialized", ""
}

// coursesIn は文字列に含まれるコースを返す。コース名がなく学科名がある場合は学科のコースを返す。
func coursesIn(s string) []string {
	s = toHalfWidth(s)
	var out []string
	add := func(cs ...string) {
		for _, c := range cs {
			if !slices.Contains(out, c) {
				out = append(out, c)
			}
		}
	}
	if strings.Contains(s, "情報システム") {
		add(courseInformationSystem)
	}
	if strings.Contains(s, "高度ICT") {
		add(courseAdvancedICT)
	}
	if strings.Contains(s, "情報デザイン") || strings.Contains(s, "デザイン") && !strings.Contains(s, "メディアデザイン") {
		add(courseInformationDesign)
	}
	if strings.Contains(strings.ReplaceAll(s, "複雑系知能学科", ""), "複雑系") {
		add(courseComplexSystem)
	}
	if strings.Contains(s, "知能システム") {
		add(courseIntelligentSystem)
	}
	if len(out) > 0 {
		return out
	}
	if strings.Contains(s, "情報アーキテクチャ学科") {
		add(courseInformationSystem, courseInformationDesign, courseAdvancedICT)
	}
	if strings.Contains(s, "複雑系知能学科") {
		add(courseComplexSystem, courseIntelligentSystem)
	}
	if strings.Contains(s, "全コース") || strings.Contains(s, "学部共通") {
		add(allCourses...)
	}
	return out
}

func requirementType(s string) string {
	switch {
	case strings.Contains(s, "選択必修"):
		return "OptionalRequired"
	case strings.Contains(s, "必修"):
		return "Required"
	case strings.Contains(s, "選択"):
		return "Optional"
	}
	return ""
}

// Requirements は対象コース・領域と科目群・科目区分から、コースごとの必修・選択を決める。
// 科目群・科目区分の各行（例: 情報システムコース専門必修）を優先し、決まらない場合は対象コース・領域を使う。
// 大学院の科目はコースを持たないため空にする。
func Requirements(targetCourses, classifications string, graduate bool) []model.SubjectRequirement {
	if graduate {
		return nil
	}
	types := map[string]string{}
	var order []string
	set := func(course, typ string) {
		if _, ok := types[course]; !ok {
			order = append(order, course)
			types[course] = typ
		}
	}
	for _, line := range lines(classifications) {
		typ := requirementType(line)
		if typ == "" {
			continue
		}
		// 例: 情報アーキテクチャ学科専門必修（高度ICT）は括弧内のコースのみが対象
		target := line
		if i := strings.IndexAny(line, "（("); i >= 0 {
			if cs := coursesIn(line[i:]); len(cs) > 0 {
				target = line[i:]
			}
		}
		for _, c := range coursesIn(target) {
			set(c, typ)
		}
	}
	if len(types) == 0 {
		typ := "Optional"
		for _, line := range lines(targetCourses) {
			if t := requirementType(line); t != "" {
				typ = t
				break
			}
		}
		if t := requirementType(classifications); t != "" {
			typ = t
		}
		var courses []string
		for _, line := range lines(targetCourses) {
			if strings.HasPrefix(line, "※") || strings.Contains(line, "受け入れる") {
				break
			}
			for _, c := range coursesIn(line) {
				if !slices.Contains(courses, c) {
					courses = append(courses, c)
				}
			}
		}
		for _, c := range courses {
			set(c, typ)
		}
	}
	out := make([]model.SubjectRequirement, 0, len(order))
	for _, c := range order {
		out = append(out, model.SubjectRequirement{Course: c, RequirementType: types[c]})
	}
	return out
}

var (
	// 授業名末尾の対象学年・クラス。例: 1-EFGH, 2-GHIの半分, 1～4, 1,2, M1,2, 1)1～4
	nameSuffixPattern = regexp.MustCompile(`(?:\d\))?(M?)(\d)(?:[～~](\d)|((?:,\d)+))?(?:-([A-L]+))?[^0-9A-Za-z]*$`)
	gradeYearPattern  = regexp.MustCompile(`(\d)年`)
	// 授業名末尾の補足（例: 「 夏期集中」「(9月入学者ほか)」）
	trailingNotePattern = regexp.MustCompile(`(\s*(夏期集中|冬期集中|\([^()]*\)))+$`)
)

// EligibleAttributes は授業名の末尾（例: 1-EFGH, 2-GHI, 1～4）と配当年次から対象学年・クラスを決める。
func EligibleAttributes(name, grades string, graduate bool) ([]model.SubjectEligibleAttribute, string) {
	prefix := "B"
	if graduate {
		prefix = "M"
	}
	suffix := toHalfWidth(strings.TrimSpace(name))
	// 「　夏期集中」「（9月入学者ほか）」などの補足を取り除く
	suffix = trailingNotePattern.ReplaceAllString(suffix, "")
	var years []int
	var classes []string
	if m := nameSuffixPattern.FindStringSubmatch(suffix); m != nil && m[0] != "" {
		if m[1] == "M" {
			prefix = "M"
		}
		from, _ := strconv.Atoi(m[2])
		years = append(years, from)
		switch {
		case m[3] != "":
			to, _ := strconv.Atoi(m[3])
			for y := from + 1; y <= to; y++ {
				years = append(years, y)
			}
		case m[4] != "":
			for _, n := range strings.Split(strings.TrimPrefix(m[4], ","), ",") {
				y, _ := strconv.Atoi(n)
				years = append(years, y)
			}
		}
		for _, c := range m[5] {
			classes = append(classes, string(c))
		}
	} else {
		for _, m := range gradeYearPattern.FindAllStringSubmatch(toHalfWidth(grades), -1) {
			y, _ := strconv.Atoi(m[1])
			if !slices.Contains(years, y) {
				years = append(years, y)
			}
		}
	}
	slices.Sort(years)
	var out []model.SubjectEligibleAttribute
	var warn string
	for _, y := range years {
		grade := fmt.Sprintf("%s%d", prefix, y)
		if !validGrade(grade) {
			warn = fmt.Sprintf("対象学年を解釈できません: %q", grade)
			continue
		}
		if len(classes) == 0 {
			out = append(out, model.SubjectEligibleAttribute{Grade: grade})
			continue
		}
		for _, c := range classes {
			class := c
			out = append(out, model.SubjectEligibleAttribute{Grade: grade, Class: &class})
		}
	}
	if len(out) == 0 && warn == "" {
		warn = fmt.Sprintf("対象学年を解釈できません: 授業名 %q, 配当年次 %q", name, grades)
	}
	return out, warn
}

func validGrade(g string) bool {
	switch g {
	case "B1", "B2", "B3", "B4", "M1", "M2", "D1", "D2", "D3":
		return true
	}
	return false
}

// isGraduate は大学院（研究科）の科目かを判定する。
func isGraduate(s model.Syllabus) bool {
	return strings.Contains(s.Classifications, "研究科") ||
		strings.Contains(s.TargetCourses, "領域") ||
		strings.Contains(s.DiplomaPolicies, "大学院")
}

// FacultyNames は担当教員名を改行で分割し、◎ の付いた教員を主担当とする。
func FacultyNames(s string) []FacultyName {
	var out []FacultyName
	for _, line := range lines(s) {
		primary := strings.HasPrefix(line, "◎")
		name := strings.TrimSpace(strings.TrimPrefix(line, "◎"))
		if name == "" {
			continue
		}
		out = append(out, FacultyName{Name: name, IsPrimary: primary})
	}
	return out
}

// normalizeName は突合のため空白（全角を含む）を取り除く。
func normalizeName(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func lines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// toHalfWidth は全角英数字を半角にする。
func toHalfWidth(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '０' && r <= '９', r >= 'Ａ' && r <= 'Ｚ', r >= 'ａ' && r <= 'ｚ':
			return r - '０' + '0'
		case r == '（':
			return '('
		case r == '）':
			return ')'
		case r == '　':
			return ' '
		}
		return r
	}, s)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
