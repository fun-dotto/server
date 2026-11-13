package subjectbuild

import (
	"reflect"
	"testing"

	"github.com/fun-dotto/server/internal/shared/model"
)

func attrs(as []model.SubjectEligibleAttribute) []string {
	var out []string
	for _, a := range as {
		s := a.Grade
		if a.Class != nil {
			s += *a.Class
		}
		out = append(out, s)
	}
	return out
}

func TestEligibleAttributes(t *testing.T) {
	tests := []struct {
		name, grades string
		graduate     bool
		want         []string
	}{
		{"認知科学1～4", "1年", false, []string{"B1", "B2", "B3", "B4"}},
		{"情報表現基礎 I1-EF", "1年", false, []string{"B1E", "B1F"}},
		{"情報処理演習I2-GHIの半分", "2年", false, []string{"B2G", "B2H", "B2I"}},
		{"線形代数学I1,2", "1年", false, []string{"B1", "B2"}},
		{"芸術論1～4　夏期集中", "1年", false, []string{"B1", "B2", "B3", "B4"}},
		{"グローバル共創プログラム (1)1～4", "1年", false, []string{"B1", "B2", "B3", "B4"}},
		{"IT・ビジネススキル4-高度ICT", "4年", false, []string{"B4"}},
		{"情報ネットワーク特論ⅡM1,2", "1年", true, []string{"M1", "M2"}},
		{"課題研究ⅡM1（9月入学者ほか）", "1年", true, []string{"M1"}},
		{"システム情報科学特別研究", "1年\n2年", true, []string{"M1", "M2"}},
	}
	for _, tt := range tests {
		got, _ := EligibleAttributes(tt.name, tt.grades, tt.graduate)
		if !reflect.DeepEqual(attrs(got), tt.want) {
			t.Errorf("EligibleAttributes(%q) = %v, want %v", tt.name, attrs(got), tt.want)
		}
	}
}

func TestSemesterAndCredit(t *testing.T) {
	if got := Semester("前期", "a"); got != "H1" {
		t.Errorf("got %q", got)
	}
	if got := Semester("3Q", "a"); got != "Q3" {
		t.Errorf("got %q", got)
	}
	if got := Semester("前期", "芸術論1～4　夏期集中"); got != "SummerIntensive" {
		t.Errorf("got %q", got)
	}
	if got := Semester("", "a"); got != "" {
		t.Errorf("got %q", got)
	}
	if n, ok := Credit("2単位"); !ok || n != 2 {
		t.Errorf("got %d %v", n, ok)
	}
	if n, ok := Credit("0単位～6単位"); !ok || n != 6 {
		t.Errorf("got %d %v", n, ok)
	}
}

func TestClassification(t *testing.T) {
	tests := []struct{ in, c, cat string }{
		{"教養基礎科目群（人間形成）", "Cultural", "Human"},
		{"教養基礎科目群（社会参加）", "Cultural", "Society"},
		{"コミュニケーション科目群", "Cultural", "Communication"},
		{"教養基礎科目群（その他）", "Cultural", ""},
		{"研究指導科目", "ResearchInstruction", ""},
		{"情報デザインコース専門必修", "Specialized", ""},
	}
	for _, tt := range tests {
		c, cat := Classification(tt.in)
		if c != tt.c || cat != tt.cat {
			t.Errorf("Classification(%q) = %q, %q", tt.in, c, cat)
		}
	}
}

func reqs(rs []model.SubjectRequirement) map[string]string {
	out := map[string]string{}
	for _, r := range rs {
		out[r.Course] = r.RequirementType
	}
	return out
}

func TestRequirements(t *testing.T) {
	tests := []struct {
		target, classifications string
		want                    map[string]string
	}{
		{"情報システム\n高度ICT", "情報システムコース専門必修\n高度ICTコース専門必修",
			map[string]string{"InformationSystem": "Required", "AdvancedICT": "Required"}},
		{"高度ICT\n情報システム", "高度ICTコース専門必修\n情報システムコース専門選択",
			map[string]string{"AdvancedICT": "Required", "InformationSystem": "Optional"}},
		{"情報システム\n高度ICT", "情報アーキテクチャ学科専門必修（高度ICT）\n情報アーキテクチャ学科専門選択（情報システム）",
			map[string]string{"AdvancedICT": "Required", "InformationSystem": "Optional"}},
		{"複雑系", "複雑系コース専門選択必修", map[string]string{"ComplexSystem": "OptionalRequired"}},
		{"全コース\n学部共通必修", "コミュニケーション科目群", map[string]string{
			"InformationSystem": "Required", "InformationDesign": "Required", "AdvancedICT": "Required",
			"ComplexSystem": "Required", "IntelligentSystem": "Required"}},
		{"知能システム", "専門科目", map[string]string{"IntelligentSystem": "Optional"}},
	}
	for _, tt := range tests {
		if got := reqs(Requirements(tt.target, tt.classifications, false)); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Requirements(%q, %q) = %v, want %v", tt.target, tt.classifications, got, tt.want)
		}
	}
	if got := Requirements("全領域", "研究科共通科目", true); len(got) != 0 {
		t.Errorf("graduate: got %v", got)
	}
}

func TestResolveFaculties(t *testing.T) {
	a := model.Faculty{Name: "花田 光彦"}
	a.ID[0] = 1
	r := Result{FacultyNames: FacultyNames("◎花田　光彦\n新任（仮）")}
	missing := r.ResolveFaculties([]model.Faculty{a})
	if len(r.Subject.Faculties) != 1 || !r.Subject.Faculties[0].IsPrimary || r.Subject.Faculties[0].FacultyID != a.ID {
		t.Errorf("faculties = %+v", r.Subject.Faculties)
	}
	if !reflect.DeepEqual(missing, []string{"新任（仮）"}) {
		t.Errorf("missing = %v", missing)
	}
}
