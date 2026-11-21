package service

import (
	"testing"

	d "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/stretchr/testify/assert"
)

func TestParseSemester(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		want   d.CourseSemester
		wantOK bool
	}{
		{name: "前期", in: "前期", want: d.CourseSemesterH1, wantOK: true},
		{name: "後期", in: "後期", want: d.CourseSemesterH2, wantOK: true},
		{name: "クォーター", in: "3Q", want: d.CourseSemesterQ3, wantOK: true},
		{name: "全角クォーター", in: "２Ｑ", want: d.CourseSemesterQ2, wantOK: true},
		{name: "通年", in: "通年", want: d.CourseSemesterAllYear, wantOK: true},
		{name: "前期集中は夏季集中", in: "前期集中", want: d.CourseSemesterSummerIntensive, wantOK: true},
		{name: "冬季集中", in: "冬季集中", want: d.CourseSemesterWinterIntensive, wantOK: true},
		{name: "前期・後期の併記は通年", in: "前期・後期", want: d.CourseSemesterAllYear, wantOK: true},
		{name: "1Q・2Q は前期", in: "1Q・2Q", want: d.CourseSemesterH1, wantOK: true},
		{name: "3Q・4Q は後期", in: "3Q・4Q", want: d.CourseSemesterH2, wantOK: true},
		{name: "空欄", in: "", want: "", wantOK: false},
		{name: "不明", in: "その他", want: "", wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseSemester(tt.in)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOK, ok)
		})
	}
}

func TestParseClassification(t *testing.T) {
	tests := []struct {
		name         string
		in           string
		wantClass    d.SubjectClassification
		wantCategory string
	}{
		{name: "専門科目", in: "専門科目 必修", wantClass: d.SubjectClassificationSpecialized, wantCategory: ""},
		{name: "教養科目のカテゴリ", in: "教養基礎科目群 社会と文化", wantClass: d.SubjectClassificationCultural, wantCategory: string(d.CulturalSubjectCategorySociety)},
		{name: "コミュニケーション", in: "教養科目 コミュニケーション", wantClass: d.SubjectClassificationCultural, wantCategory: string(d.CulturalSubjectCategoryCommunication)},
		{name: "教養でカテゴリ不明", in: "教養科目", wantClass: d.SubjectClassificationCultural, wantCategory: ""},
		{name: "研究指導", in: "研究指導科目", wantClass: d.SubjectClassificationResearchInstruction, wantCategory: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			class, category := parseClassification(tt.in)
			assert.Equal(t, tt.wantClass, class)
			assert.Equal(t, tt.wantCategory, category)
		})
	}
}

func TestParseFacultyNames(t *testing.T) {
	got := parseFacultyNames("◎函館　太郎\r\n未来 花子\n\n未来花子")
	want := []syllabusFacultyName{
		{Name: "函館太郎", IsPrimary: true},
		{Name: "未来花子", IsPrimary: false},
	}
	assert.Equal(t, want, got)
}

func ptr(s string) *string { return &s }

func TestParseEligibleAttributes(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		grades  string
		want    []model.SubjectEligibleAttribute
	}{
		{name: "授業名末尾のクラス指定を優先", subject: "解析学Ⅰ 2-GHI", grades: "1年", want: []model.SubjectEligibleAttribute{
			{Grade: "B2", Class: ptr("G")}, {Grade: "B2", Class: ptr("H")}, {Grade: "B2", Class: ptr("I")},
		}},
		{name: "全角のクラス指定", subject: "解析学Ⅰ　１－ＥＦＧＨ", grades: "", want: []model.SubjectEligibleAttribute{
			{Grade: "B1", Class: ptr("E")}, {Grade: "B1", Class: ptr("F")}, {Grade: "B1", Class: ptr("G")}, {Grade: "B1", Class: ptr("H")},
		}},
		{name: "授業名末尾の学年範囲", subject: "認知科学1～4", grades: "", want: []model.SubjectEligibleAttribute{
			{Grade: "B1"}, {Grade: "B2"}, {Grade: "B3"}, {Grade: "B4"},
		}},
		{name: "授業名の半角ハイフンは範囲とみなさない", subject: "情報処理演習Ⅱ 2-3", grades: "2年", want: []model.SubjectEligibleAttribute{{Grade: "B2"}}},
		{name: "配当年次の単一学年", subject: "認知科学", grades: "3年", want: []model.SubjectEligibleAttribute{{Grade: "B3"}}},
		{name: "配当年次の列挙", subject: "認知科学", grades: "2,3年", want: []model.SubjectEligibleAttribute{{Grade: "B2"}, {Grade: "B3"}}},
		{name: "年の付いた範囲", subject: "認知科学", grades: "1年～4年", want: []model.SubjectEligibleAttribute{{Grade: "B1"}, {Grade: "B2"}, {Grade: "B3"}, {Grade: "B4"}}},
		{name: "以上", subject: "認知科学", grades: "2年以上", want: []model.SubjectEligibleAttribute{{Grade: "B2"}, {Grade: "B3"}, {Grade: "B4"}}},
		{name: "大学院修士", subject: "特論", grades: "修士1～2年", want: []model.SubjectEligibleAttribute{{Grade: "M1"}, {Grade: "M2"}}},
		{name: "博士（前期）は修士", subject: "特論", grades: "博士（前期）1年", want: []model.SubjectEligibleAttribute{{Grade: "M1"}}},
		{name: "博士", subject: "特論", grades: "博士1年", want: []model.SubjectEligibleAttribute{{Grade: "D1"}}},
		{name: "解釈できない", subject: "特論", grades: "", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseEligibleAttributes(tt.subject, tt.grades))
		})
	}
}

func TestParseRequirements(t *testing.T) {
	req := func(c d.CourseType, r d.SubjectRequirementType) model.SubjectRequirement {
		return model.SubjectRequirement{Course: string(c), RequirementType: string(r)}
	}
	allCourses := func(r d.SubjectRequirementType) []model.SubjectRequirement {
		return []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, r),
			req(d.CourseTypeInformationDesign, r),
			req(d.CourseTypeAdvancedICT, r),
			req(d.CourseTypeComplexSystem, r),
			req(d.CourseTypeIntelligentSystem, r),
		}
	}
	tests := []struct {
		name            string
		targetCourses   string
		classifications string
		want            []model.SubjectRequirement
	}{
		{name: "コース指定なしは対象コース全体に適用", targetCourses: "情報システムコース、情報デザインコース", classifications: "専門科目 必修", want: []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeRequired),
		}},
		{name: "コースごとの指定と選択必修", targetCourses: "情報システムコース、高度ICTコース", classifications: "情報システムコース 選択必修\n高度ICTコース 必修", want: []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeOptionalRequired),
			req(d.CourseTypeAdvancedICT, d.SubjectRequirementTypeRequired),
		}},
		{name: "記載のない対象コースは選択", targetCourses: "複雑系コース、知能システムコース", classifications: "知能システムコース 必修", want: []model.SubjectRequirement{
			req(d.CourseTypeIntelligentSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeComplexSystem, d.SubjectRequirementTypeOptional),
		}},
		{name: "対象コースなし", targetCourses: "", classifications: "教養科目 選択", want: []model.SubjectRequirement{}},
		{name: "1 つの区切りにコースと区分が交互に並ぶ", targetCourses: "情報システムコース、情報デザインコース", classifications: "情報システムコース 必修 情報デザインコース 選択", want: []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeOptional),
		}},
		{name: "複数コースの後ろにまとめた区分", targetCourses: "情報システムコース、情報デザインコース", classifications: "情報システムコース・情報デザインコース 必修", want: []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeRequired),
		}},
		{name: "区分がコース名より前", targetCourses: "情報システムコース、情報デザインコース", classifications: "必修 情報システムコース 選択 情報デザインコース", want: []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeOptional),
		}},
		{name: "科目区分欄を対象コース欄より優先", targetCourses: "全コース\n学部共通選択", classifications: "専門科目 必修", want: allCourses(d.SubjectRequirementTypeRequired)},
		{name: "全コースと対象コース欄の必修・選択", targetCourses: "全コース\n学部共通選択", classifications: "共通教養科目", want: allCourses(d.SubjectRequirementTypeOptional)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseRequirements(tt.targetCourses, tt.classifications))
		})
	}
}
