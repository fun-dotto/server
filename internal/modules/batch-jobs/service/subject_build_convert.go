package service

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"golang.org/x/text/width"
)

// SyllabusFacultyName はシラバスの担当教員名を 1 人分に分解したもの。
type SyllabusFacultyName struct {
	Name      string
	IsPrimary bool
}

// StructuredSubject はシラバスを構造化した科目。担当教員は名前のままで、faculties との突合は呼び出し側で行う。
type StructuredSubject struct {
	Name                    string
	Year                    int
	Semester                academicdomain.CourseSemester
	Credit                  int
	Classification          academicdomain.SubjectClassification
	CulturalSubjectCategory string
	FacultyNames            []SyllabusFacultyName
	EligibleAttributes      []model.SubjectEligibleAttribute
	Requirements            []model.SubjectRequirement
}

// StructureSyllabus はシラバスの文字列を科目の構造化データに変換する。
// 開講期が解釈できない場合は semester が必須の enum のため ok=false を返す。
func StructureSyllabus(s model.Syllabus) (StructuredSubject, bool) {
	semester, ok := ParseSemester(s.Term)
	if !ok {
		return StructuredSubject{}, false
	}
	classification, cultural := ParseClassification(s.Classifications)
	return StructuredSubject{
		Name:                    strings.TrimSpace(s.Name),
		Year:                    s.Year,
		Semester:                semester,
		Credit:                  model.ParseSyllabusCredit(s.Credit),
		Classification:          classification,
		CulturalSubjectCategory: cultural,
		FacultyNames:            ParseFacultyNames(s.FacultyNames),
		EligibleAttributes:      ParseEligibleAttributes(s.Name, s.Grades),
		Requirements:            ParseRequirements(s.TargetCourses, s.Classifications),
	}, true
}

// dashVariants はハイフンとして書かれうる文字。width.Fold では半角にならない。
var dashVariants = strings.NewReplacer("ー", "-", "−", "-", "‐", "-", "–", "-", "―", "-")

// normalizeWidth は全角の英数字・記号を半角に揃え、ハイフンの異体字を `-` にする。
func normalizeWidth(s string) string {
	return dashVariants.Replace(width.Fold.String(s))
}

// ParseSemester は開講期（`前期`/`後期`/`1Q`〜`4Q`/`通年`/集中）を enum に変換する。
func ParseSemester(term string) (academicdomain.CourseSemester, bool) {
	t := strings.ToUpper(normalizeWidth(strings.Join(strings.Fields(term), "")))
	isIntensive := strings.Contains(t, "集中")
	hasQ := func(n string) bool { return strings.Contains(t, n+"Q") }
	switch {
	case t == "":
		return "", false
	case strings.Contains(t, "通年"):
		return academicdomain.CourseSemesterAllYear, true
	// 開講期が併記されている場合は、併記した期間全体を表す値にする
	case !isIntensive && strings.Contains(t, "前期") && strings.Contains(t, "後期"),
		hasQ("1") && hasQ("2") && hasQ("3") && hasQ("4"):
		return academicdomain.CourseSemesterAllYear, true
	case hasQ("1") && hasQ("2") && !hasQ("3") && !hasQ("4"):
		return academicdomain.CourseSemesterH1, true
	case hasQ("3") && hasQ("4") && !hasQ("1") && !hasQ("2"):
		return academicdomain.CourseSemesterH2, true
	// 集中講義は「前期集中」のように前期・後期と併記されうるため先に判定する
	case strings.Contains(t, "集中") && (strings.Contains(t, "夏") || strings.Contains(t, "前期")):
		return academicdomain.CourseSemesterSummerIntensive, true
	case strings.Contains(t, "集中") && (strings.Contains(t, "冬") || strings.Contains(t, "後期")):
		return academicdomain.CourseSemesterWinterIntensive, true
	case strings.Contains(t, "1Q"):
		return academicdomain.CourseSemesterQ1, true
	case strings.Contains(t, "2Q"):
		return academicdomain.CourseSemesterQ2, true
	case strings.Contains(t, "3Q"):
		return academicdomain.CourseSemesterQ3, true
	case strings.Contains(t, "4Q"):
		return academicdomain.CourseSemesterQ4, true
	case strings.Contains(t, "前期"):
		return academicdomain.CourseSemesterH1, true
	case strings.Contains(t, "後期"):
		return academicdomain.CourseSemesterH2, true
	}
	return "", false
}

// ParseClassification は科目群・科目区分から科目区分と教養科目カテゴリを判定する。
// 教養・研究指導の記載がなければ専門科目とみなす。教養科目以外のカテゴリは空文字。
func ParseClassification(classifications string) (academicdomain.SubjectClassification, string) {
	switch {
	case strings.Contains(classifications, "研究指導"):
		return academicdomain.SubjectClassificationResearchInstruction, ""
	case strings.Contains(classifications, "教養"):
		return academicdomain.SubjectClassificationCultural, string(parseCulturalCategory(classifications))
	}
	return academicdomain.SubjectClassificationSpecialized, ""
}

func parseCulturalCategory(s string) academicdomain.CulturalSubjectCategory {
	// 「人間」を含む語（例: 人間形成）より先に、語の重なりがない区分から判定する
	switch {
	case strings.Contains(s, "コミュニケーション"):
		return academicdomain.CulturalSubjectCategoryCommunication
	case strings.Contains(s, "健康"):
		return academicdomain.CulturalSubjectCategoryHealth
	case strings.Contains(s, "社会"):
		return academicdomain.CulturalSubjectCategorySociety
	case strings.Contains(s, "科学"):
		return academicdomain.CulturalSubjectCategoryScience
	case strings.Contains(s, "人間"):
		return academicdomain.CulturalSubjectCategoryHuman
	}
	return ""
}

// ParseFacultyNames は改行区切りの担当教員名を分解する。`◎` の付いた教員が主担当。
func ParseFacultyNames(s string) []SyllabusFacultyName {
	var result []SyllabusFacultyName
	seen := map[string]bool{}
	for _, line := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		isPrimary := strings.Contains(line, "◎")
		name := NormalizeFacultyName(strings.ReplaceAll(line, "◎", ""))
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		result = append(result, SyllabusFacultyName{Name: name, IsPrimary: isPrimary})
	}
	return result
}

// NormalizeFacultyName は教員名の突合のため空白（全角含む）を除く。シラバスと faculties で姓名間の空白の有無が揃わないため。
func NormalizeFacultyName(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "　", " ")), "")
}

var (
	// 授業名末尾のクラス指定（例: `解析学Ⅰ 1-EFGH`）
	nameGradeClassPattern = regexp.MustCompile(`([1-4])\s*-\s*([A-L]+)\s*$`)
	// 授業名末尾の学年範囲（例: `認知科学1～4`）。`2-3` のような半角ハイフンは授業の番号と区別できないため範囲とみなさない
	nameGradeRangePattern = regexp.MustCompile(`([1-4])\s*[～~〜]\s*([1-4])\s*$`)
	gradeRangePattern     = regexp.MustCompile(`([1-4])\s*年?\s*[～~〜-]\s*([1-4])`)
	gradeAbovePattern     = regexp.MustCompile(`([1-4])\s*年?\s*次?\s*以上`)
	gradeDigitPattern     = regexp.MustCompile(`[1-4]`)
)

// ParseEligibleAttributes は授業名の末尾と配当年次から対象の学年・クラスを求める。
// 授業名末尾のクラス指定が最も具体的なので優先し、なければ配当年次を使う。
// 大学院（修士・博士）はクラスがないため Class は nil。
func ParseEligibleAttributes(name, grades string) []model.SubjectEligibleAttribute {
	name = strings.TrimSpace(normalizeWidth(name))
	if m := nameGradeClassPattern.FindStringSubmatch(name); m != nil {
		grade := "B" + m[1]
		var result []model.SubjectEligibleAttribute
		for _, c := range m[2] {
			class := string(c)
			result = append(result, model.SubjectEligibleAttribute{Grade: grade, Class: &class})
		}
		return result
	}
	if m := nameGradeRangePattern.FindStringSubmatch(name); m != nil {
		return gradeRange("B", m[1], m[2])
	}

	g := normalizeWidth(grades)
	prefix := "B"
	maxGrade := "4"
	// この大学の修士課程の正式名は「博士（前期）課程」なので、前期・後期を先に判定する
	gp := strings.NewReplacer("（", "(", "）", ")").Replace(g)
	switch {
	case strings.Contains(gp, "博士(前期)"), strings.Contains(gp, "博士前期"), strings.Contains(g, "修士"), strings.Contains(g, "大学院"):
		prefix, maxGrade = "M", "2"
	case strings.Contains(g, "博士"):
		prefix, maxGrade = "D", "3"
	}
	if m := gradeRangePattern.FindStringSubmatch(g); m != nil {
		return gradeRange(prefix, m[1], m[2])
	}
	if m := gradeAbovePattern.FindStringSubmatch(g); m != nil {
		return gradeRange(prefix, m[1], maxGrade)
	}
	var result []model.SubjectEligibleAttribute
	seen := map[string]bool{}
	for _, d := range gradeDigitPattern.FindAllString(g, -1) {
		if seen[d] {
			continue
		}
		seen[d] = true
		result = append(result, model.SubjectEligibleAttribute{Grade: prefix + d})
	}
	return result
}

func gradeRange(prefix, from, to string) []model.SubjectEligibleAttribute {
	start, _ := strconv.Atoi(from)
	end, _ := strconv.Atoi(to)
	var result []model.SubjectEligibleAttribute
	for i := start; i <= end; i++ {
		result = append(result, model.SubjectEligibleAttribute{Grade: prefix + strconv.Itoa(i)})
	}
	return result
}

var courseKeywords = []struct {
	keyword string
	course  academicdomain.CourseType
}{
	{"情報システム", academicdomain.CourseTypeInformationSystem},
	{"情報デザイン", academicdomain.CourseTypeInformationDesign},
	{"高度ICT", academicdomain.CourseTypeAdvancedICT},
	{"複雑系", academicdomain.CourseTypeComplexSystem},
	{"知能システム", academicdomain.CourseTypeIntelligentSystem},
}

func findCourses(s string) []academicdomain.CourseType {
	s = strings.ToUpper(s)
	var result []academicdomain.CourseType
	if strings.Contains(s, "全コース") {
		for _, k := range courseKeywords {
			result = append(result, k.course)
		}
		return result
	}
	for _, k := range courseKeywords {
		if strings.Contains(s, k.keyword) {
			result = append(result, k.course)
		}
	}
	return result
}

func findRequirementType(s string) (academicdomain.SubjectRequirementType, bool) {
	// 「選択必修」は「必修」「選択」の両方を含むため最初に判定する
	switch {
	case strings.Contains(s, "選択必修"):
		return academicdomain.SubjectRequirementTypeOptionalRequired, true
	case strings.Contains(s, "必修"):
		return academicdomain.SubjectRequirementTypeRequired, true
	case strings.Contains(s, "選択"):
		return academicdomain.SubjectRequirementTypeOptional, true
	}
	return "", false
}

var segmentSeparator = regexp.MustCompile(`[\n、,，/／;；]+`)

var (
	coursePattern      = regexp.MustCompile(`全コース|情報システム|情報デザイン|高度ICT|複雑系|知能システム`)
	requirementPattern = regexp.MustCompile(`選択必修|必修|選択`)
)

type courseRequirement struct {
	courses     []academicdomain.CourseType
	requirement academicdomain.SubjectRequirementType
}

// parseSegment は 1 つの区切りの中のコース名と必修・選択を、出てくる順に組にする。
// 「A B 必修 C 選択」のようにコース名が先なら、区分はその前に並んだコースすべてに付く。
// 「必修 A 選択 B」のように区分が先なら、区分はその後に並んだコースに付く。
// 組になるコースがない区分は、targets（対象コース全体）に付く。
func parseSegment(seg string, targets []academicdomain.CourseType) []courseRequirement {
	upper := strings.ToUpper(seg)
	type token struct {
		pos     int
		courses []academicdomain.CourseType
		req     academicdomain.SubjectRequirementType
	}
	var tokens []token
	for _, m := range coursePattern.FindAllStringIndex(upper, -1) {
		tokens = append(tokens, token{pos: m[0], courses: findCourses(upper[m[0]:m[1]])})
	}
	for _, m := range requirementPattern.FindAllStringIndex(seg, -1) {
		r, _ := findRequirementType(seg[m[0]:m[1]])
		tokens = append(tokens, token{pos: m[0], req: r})
	}
	slices.SortFunc(tokens, func(a, b token) int { return a.pos - b.pos })
	if len(tokens) == 0 {
		return nil
	}

	reqFirst := tokens[0].req != ""
	var result []courseRequirement
	var pending []academicdomain.CourseType
	var current academicdomain.SubjectRequirementType
	flush := func() {
		if current == "" {
			return
		}
		courses := pending
		if len(courses) == 0 {
			courses = targets
		}
		result = append(result, courseRequirement{courses: courses, requirement: current})
		pending, current = nil, ""
	}
	for _, t := range tokens {
		switch {
		case t.req == "":
			pending = append(pending, t.courses...)
		case reqFirst:
			// 区分が先の書き方では、新しい区分が出たところで直前の区分とコースの組を確定する
			flush()
			current = t.req
		default:
			current = t.req
			flush()
		}
	}
	if reqFirst {
		flush()
	}
	return result
}

// ParseRequirements は対象コースと科目群・科目区分から、コースごとの必修・選択を求める。
// 区切りごとにコース名があればそのコースに、なければ対象コース全体に必修・選択を割り当てる。「全コース」は 5 コースすべてとみなす。
// 科目区分欄の記載を対象コース欄（「学部共通選択」など）より優先し、同じ欄の中では先に現れた割り当てを優先する。
// どの区切りにも必修・選択の記載がない対象コースは選択とみなす。
func ParseRequirements(targetCourses, classifications string) []model.SubjectRequirement {
	targets := findCourses(targetCourses)
	assigned := map[academicdomain.CourseType]academicdomain.SubjectRequirementType{}
	var order []academicdomain.CourseType
	assign := func(c academicdomain.CourseType, r academicdomain.SubjectRequirementType) {
		if _, ok := assigned[c]; ok {
			return
		}
		assigned[c] = r
		order = append(order, c)
	}

	for _, field := range []string{classifications, targetCourses} {
		for _, seg := range segmentSeparator.Split(field, -1) {
			for _, cr := range parseSegment(seg, targets) {
				for _, c := range cr.courses {
					assign(c, cr.requirement)
				}
			}
		}
	}
	for _, c := range targets {
		assign(c, academicdomain.SubjectRequirementTypeOptional)
	}

	result := make([]model.SubjectRequirement, 0, len(order))
	for _, c := range order {
		result = append(result, model.SubjectRequirement{Course: string(c), RequirementType: string(assigned[c])})
	}
	return result
}
