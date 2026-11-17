package service

import (
	"regexp"
	"strconv"
	"strings"

	academicdomain "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
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

var fullWidthDigits = strings.NewReplacer("０", "0", "１", "1", "２", "2", "３", "3", "４", "4", "５", "5", "６", "6", "７", "7", "８", "8", "９", "9", "Ｑ", "Q", "ｑ", "Q")

// ParseSemester は開講期（`前期`/`後期`/`1Q`〜`4Q`/`通年`/集中）を enum に変換する。
func ParseSemester(term string) (academicdomain.CourseSemester, bool) {
	t := strings.ToUpper(fullWidthDigits.Replace(strings.Join(strings.Fields(term), "")))
	switch {
	case t == "":
		return "", false
	case strings.Contains(t, "通年"):
		return academicdomain.CourseSemesterAllYear, true
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
	// 授業名末尾の学年範囲（例: `xxx 1～4`）
	nameGradeRangePattern = regexp.MustCompile(`\s([1-4])\s*[～~〜-]\s*([1-4])\s*$`)
	gradeRangePattern     = regexp.MustCompile(`([1-4])\s*年?\s*[～~〜-]\s*([1-4])`)
	gradeAbovePattern     = regexp.MustCompile(`([1-4])\s*年?\s*次?\s*以上`)
	gradeDigitPattern     = regexp.MustCompile(`[1-4]`)
)

// ParseEligibleAttributes は授業名の末尾と配当年次から対象の学年・クラスを求める。
// 授業名末尾のクラス指定が最も具体的なので優先し、なければ配当年次を使う。
// 大学院（修士・博士）はクラスがないため Class は nil。
func ParseEligibleAttributes(name, grades string) []model.SubjectEligibleAttribute {
	name = fullWidthDigits.Replace(strings.TrimSpace(name))
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

	g := fullWidthDigits.Replace(grades)
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

// ParseRequirements は対象コースと科目群・科目区分から、コースごとの必修・選択を求める。
// 対象コースと科目区分の各区切りにコース名があればそのコースに、なければ対象コース全体に必修・選択を割り当てる。
// 「全コース」は 5 コースすべてとみなす。
// 先に現れた割り当てを優先し、どの区切りにも必修・選択の記載がないコースは選択とみなす。
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

	// 対象コース欄にも「学部共通選択」のように必修・選択が書かれることがあるため、両方の区切りを見る
	for _, seg := range segmentSeparator.Split(targetCourses+"\n"+classifications, -1) {
		r, ok := findRequirementType(seg)
		if !ok {
			continue
		}
		courses := findCourses(seg)
		if len(courses) == 0 {
			courses = targets
		}
		for _, c := range courses {
			assign(c, r)
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
