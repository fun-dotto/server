package service

import (
	"reflect"
	"testing"

	d "github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func TestParseSemester(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want d.CourseSemester
		ok   bool
	}{
		{"前期", "前期", d.CourseSemesterH1, true},
		{"後期", "後期", d.CourseSemesterH2, true},
		{"クォーター", "3Q", d.CourseSemesterQ3, true},
		{"全角クォーター", "２Ｑ", d.CourseSemesterQ2, true},
		{"通年", "通年", d.CourseSemesterAllYear, true},
		{"前期集中は夏季集中", "前期集中", d.CourseSemesterSummerIntensive, true},
		{"冬季集中", "冬季集中", d.CourseSemesterWinterIntensive, true},
		{"前期・後期の併記は通年", "前期・後期", d.CourseSemesterAllYear, true},
		{"1Q・2Q は前期", "1Q・2Q", d.CourseSemesterH1, true},
		{"3Q・4Q は後期", "3Q・4Q", d.CourseSemesterH2, true},
		{"空欄", "", "", false},
		{"不明", "その他", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseSemester(tt.in)
			if got != tt.want || ok != tt.ok {
				t.Errorf("ParseSemester(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
			}
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
		{"専門科目", "専門科目 必修", d.SubjectClassificationSpecialized, ""},
		{"教養科目のカテゴリ", "教養基礎科目群 社会と文化", d.SubjectClassificationCultural, string(d.CulturalSubjectCategorySociety)},
		{"コミュニケーション", "教養科目 コミュニケーション", d.SubjectClassificationCultural, string(d.CulturalSubjectCategoryCommunication)},
		{"教養でカテゴリ不明", "教養科目", d.SubjectClassificationCultural, ""},
		{"研究指導", "研究指導科目", d.SubjectClassificationResearchInstruction, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, cat := ParseClassification(tt.in)
			if c != tt.wantClass || cat != tt.wantCategory {
				t.Errorf("ParseClassification(%q) = %q, %q; want %q, %q", tt.in, c, cat, tt.wantClass, tt.wantCategory)
			}
		})
	}
}

func TestParseFacultyNames(t *testing.T) {
	got := ParseFacultyNames("◎函館　太郎\r\n未来 花子\n\n未来花子")
	want := []SyllabusFacultyName{
		{Name: "函館太郎", IsPrimary: true},
		{Name: "未来花子", IsPrimary: false},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseFacultyNames = %+v; want %+v", got, want)
	}
}

func ptr(s string) *string { return &s }

func TestParseEligibleAttributes(t *testing.T) {
	tests := []struct {
		name   string
		inName string
		grades string
		want   []model.SubjectEligibleAttribute
	}{
		{"授業名末尾のクラス指定を優先", "解析学Ⅰ 2-GHI", "1年", []model.SubjectEligibleAttribute{
			{Grade: "B2", Class: ptr("G")}, {Grade: "B2", Class: ptr("H")}, {Grade: "B2", Class: ptr("I")},
		}},
		{"授業名末尾の学年範囲", "英語 1～3", "", []model.SubjectEligibleAttribute{
			{Grade: "B1"}, {Grade: "B2"}, {Grade: "B3"},
		}},
		{"配当年次の単一学年", "認知科学", "3年", []model.SubjectEligibleAttribute{{Grade: "B3"}}},
		{"配当年次の列挙", "認知科学", "2,3年", []model.SubjectEligibleAttribute{{Grade: "B2"}, {Grade: "B3"}}},
		{"大学院修士", "特論", "修士1～2年", []model.SubjectEligibleAttribute{{Grade: "M1"}, {Grade: "M2"}}},
		{"博士", "特論", "博士1年", []model.SubjectEligibleAttribute{{Grade: "D1"}}},
		{"授業名の半角ハイフンは範囲とみなさない", "情報処理演習Ⅱ 2-3", "2年", []model.SubjectEligibleAttribute{{Grade: "B2"}}},
		{"空白のない授業名末尾の範囲", "認知科学1～4", "", []model.SubjectEligibleAttribute{
			{Grade: "B1"}, {Grade: "B2"}, {Grade: "B3"}, {Grade: "B4"},
		}},
		{"全角のクラス指定", "解析学Ⅰ　１－ＥＦＧＨ", "", []model.SubjectEligibleAttribute{
			{Grade: "B1", Class: ptr("E")}, {Grade: "B1", Class: ptr("F")}, {Grade: "B1", Class: ptr("G")}, {Grade: "B1", Class: ptr("H")},
		}},
		{"年の付いた範囲", "認知科学", "1年～4年", []model.SubjectEligibleAttribute{{Grade: "B1"}, {Grade: "B2"}, {Grade: "B3"}, {Grade: "B4"}}},
		{"以上", "認知科学", "2年以上", []model.SubjectEligibleAttribute{{Grade: "B2"}, {Grade: "B3"}, {Grade: "B4"}}},
		{"博士（前期）は修士", "特論", "博士（前期）1年", []model.SubjectEligibleAttribute{{Grade: "M1"}}},
		{"博士（後期）", "特論", "博士(後期)1年", []model.SubjectEligibleAttribute{{Grade: "D1"}}},
		{"解釈できない", "特論", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseEligibleAttributes(tt.inName, tt.grades)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEligibleAttributes(%q, %q) = %+v; want %+v", tt.inName, tt.grades, got, tt.want)
			}
		})
	}
}

func TestParseRequirements(t *testing.T) {
	req := func(c d.CourseType, r d.SubjectRequirementType) model.SubjectRequirement {
		return model.SubjectRequirement{Course: string(c), RequirementType: string(r)}
	}
	tests := []struct {
		name    string
		targets string
		classes string
		want    []model.SubjectRequirement
	}{
		{"コース指定なしは対象コース全体に適用", "情報システムコース、情報デザインコース", "専門科目 必修", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeRequired),
		}},
		{"コースごとの指定と選択必修", "情報システムコース、高度ICTコース", "情報システムコース 選択必修\n高度ICTコース 必修", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeOptionalRequired),
			req(d.CourseTypeAdvancedICT, d.SubjectRequirementTypeRequired),
		}},
		{"記載のない対象コースは選択", "複雑系コース、知能システムコース", "知能システムコース 必修", []model.SubjectRequirement{
			req(d.CourseTypeIntelligentSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeComplexSystem, d.SubjectRequirementTypeOptional),
		}},
		{"対象コースなし", "", "教養科目 選択", []model.SubjectRequirement{}},
		{"空白区切りで複数コース", "情報システムコース、情報デザインコース", "情報システムコース 必修 情報デザインコース 選択", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeOptional),
		}},
		{"複数コースの後ろにまとめた区分", "情報システムコース、情報デザインコース", "情報システムコース・情報デザインコース 必修", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeRequired),
		}},
		{"空白区切りの複数コースの後ろにまとめた区分", "情報システムコース、情報デザインコース", "情報システムコース 情報デザインコース 必修", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeRequired),
		}},
		{"区分がコース名より前", "情報システムコース、情報デザインコース", "必修 情報システムコース 選択 情報デザインコース", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeOptional),
		}},
		{"科目区分欄を対象コース欄より優先", "全コース\n学部共通選択", "専門科目 必修", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeAdvancedICT, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeComplexSystem, d.SubjectRequirementTypeRequired),
			req(d.CourseTypeIntelligentSystem, d.SubjectRequirementTypeRequired),
		}},
		{"全コースと対象コース欄の必修・選択", "全コース\n学部共通選択", "共通教養科目", []model.SubjectRequirement{
			req(d.CourseTypeInformationSystem, d.SubjectRequirementTypeOptional),
			req(d.CourseTypeInformationDesign, d.SubjectRequirementTypeOptional),
			req(d.CourseTypeAdvancedICT, d.SubjectRequirementTypeOptional),
			req(d.CourseTypeComplexSystem, d.SubjectRequirementTypeOptional),
			req(d.CourseTypeIntelligentSystem, d.SubjectRequirementTypeOptional),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRequirements(tt.targets, tt.classes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRequirements = %+v; want %+v", got, tt.want)
			}
		})
	}
}
