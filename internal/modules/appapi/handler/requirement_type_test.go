package handler

import (
	"strings"
	"testing"

	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

// allCourseValues は 5 コースの DB 値。courseMap を range すると順序が毎回変わるので、テストでは固定の順で並べる。
var allCourseValues = []string{"InformationSystem", "InformationDesign", "AdvancedICT", "ComplexSystem", "IntelligentSystem"}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// testUser は学年・コース・クラスを持つユーザーを作る。空文字は未入力（NULL）。
func testUser(grade, course, class string) *model.User {
	return &model.User{ID: "user1", Grade: nilIfEmpty(grade), Course: nilIfEmpty(course), Class: nilIfEmpty(class)}
}

// eligible は対象学年・クラスの行を作る。"B1-G" はクラス付き、"M1" はクラスなしの行。
func eligible(attrs ...string) []model.SubjectEligibleAttribute {
	result := make([]model.SubjectEligibleAttribute, 0, len(attrs))
	for _, a := range attrs {
		grade, class, _ := strings.Cut(a, "-")
		result = append(result, model.SubjectEligibleAttribute{Grade: grade, Class: nilIfEmpty(class)})
	}
	return result
}

func requirement(course, requirementType string) model.SubjectRequirement {
	return model.SubjectRequirement{Course: course, RequirementType: requirementType}
}

// sameForCourses は、指定したコースに同じ区分の行を作る。
func sameForCourses(requirementType string, courses ...string) []model.SubjectRequirement {
	result := make([]model.SubjectRequirement, 0, len(courses))
	for _, c := range courses {
		result = append(result, requirement(c, requirementType))
	}
	return result
}

func TestRequirementTypeFor(t *testing.T) {
	const (
		required         = appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_REQUIRED
		optional         = appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL
		optionalRequired = appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL_REQUIRED
		unspecified      = appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_UNSPECIFIED
	)

	// 実データの例: プログラミング基礎1-GH（B1-G, B1-H、全5コース必修）
	programmingBasics := &model.Subject{
		EligibleAttributes: eligible("B1-G", "B1-H"),
		Requirements:       sameForCourses("Required", allCourseValues...),
	}
	// 実データの例: 力学基礎2-GHI（B2-G, B2-H, B2-I、複雑系のみ必修。ほかのコースの行はない）
	mechanics := &model.Subject{
		EligibleAttributes: eligible("B2-G", "B2-H", "B2-I"),
		Requirements:       []model.SubjectRequirement{requirement("ComplexSystem", "Required")},
	}
	// 院生向けの科目（対象行はクラスなし）
	graduate := &model.Subject{
		EligibleAttributes: eligible("M1"),
		Requirements:       sameForCourses("Optional", allCourseValues...),
	}
	// 対象学年・クラスが空の科目（大学院の研究指導）
	noEligible := &model.Subject{
		Requirements: sameForCourses("Required", allCourseValues...),
	}

	tests := []struct {
		name    string
		user    *model.User
		subject *model.Subject
		want    appv1.SubjectRequirementType
	}{
		// 共通必修・1年生
		{"1年生・コース未設定・共通必修", testUser("B1", "", "G"), programmingBasics, required},
		{"1年生・コース選択済み・共通必修", testUser("B1", "InformationSystem", "G"), programmingBasics, required},
		{"1年生・全コースで選択", testUser("B1", "", "G"), &model.Subject{
			EligibleAttributes: eligible("B1-G"),
			Requirements:       sameForCourses("Optional", allCourseValues...),
		}, optional},
		{"1年生・コース未設定・コースによって区分が違う", testUser("B1", "", "G"), &model.Subject{
			EligibleAttributes: eligible("B1-G"),
			Requirements: append(
				sameForCourses("Required", "InformationSystem", "InformationDesign", "AdvancedICT", "ComplexSystem"),
				requirement("IntelligentSystem", "Optional"),
			),
		}, unspecified},
		{"1年生・コース未設定・行のないコースがある", testUser("B1", "", "G"), &model.Subject{
			EligibleAttributes: eligible("B1-G"),
			Requirements:       sameForCourses("Required", "InformationSystem", "InformationDesign", "AdvancedICT", "ComplexSystem"),
		}, unspecified},

		// クラス別・学年別
		{"クラスが対象外", testUser("B1", "", "A"), programmingBasics, unspecified},
		{"学年が対象外", testUser("B2", "ComplexSystem", "I"), programmingBasics, unspecified},

		// コース別
		{"コース別必修・該当コース", testUser("B2", "ComplexSystem", "I"), mechanics, required},
		{"コース別必修・行のないコース", testUser("B2", "InformationDesign", "I"), mechanics, unspecified},
		{"コース別必修・1年生は対象外", testUser("B1", "", "G"), mechanics, unspecified},
		{"選択必修", testUser("B3", "AdvancedICT", "A"), &model.Subject{
			EligibleAttributes: eligible("B3-A"),
			Requirements:       []model.SubjectRequirement{requirement("AdvancedICT", "OptionalRequired")},
		}, optionalRequired},
		{"同じコースに異なる区分の行がある", testUser("B2", "ComplexSystem", "I"), &model.Subject{
			EligibleAttributes: eligible("B2-I"),
			Requirements:       []model.SubjectRequirement{requirement("ComplexSystem", "Required"), requirement("ComplexSystem", "Optional")},
		}, unspecified},
		{"同じコースに同じ区分の行が重複している", testUser("B2", "ComplexSystem", "I"), &model.Subject{
			EligibleAttributes: eligible("B2-I"),
			Requirements:       []model.SubjectRequirement{requirement("ComplexSystem", "Required"), requirement("ComplexSystem", "Required")},
		}, required},

		// 院生
		{"院生・クラス未設定", testUser("M1", "InformationSystem", ""), graduate, optional},
		{"院生・クラスあり（学年だけで判定）", testUser("M1", "InformationSystem", "A"), graduate, optional},
		{"院生・学年が対象外", testUser("M2", "InformationSystem", ""), graduate, unspecified},

		// 対象学年・クラスが空の科目
		{"対象が空の科目は全員が対象", testUser("M1", "InformationSystem", ""), noEligible, required},

		// ユーザー不在・未入力（対象が空の科目を使い、未入力の判定だけで UNSPECIFIED になることを確かめる）
		{"ユーザーが存在しない", nil, noEligible, unspecified},
		{"すべて未入力", testUser("", "", ""), noEligible, unspecified},
		{"学年だけ未入力", testUser("", "InformationSystem", "A"), noEligible, unspecified},
		{"2年生以降でコース未入力", testUser("B2", "", "G"), noEligible, unspecified},
		{"2年生以降でコース未入力・全コースで同じ区分の科目", testUser("B2", "", "G"), &model.Subject{
			EligibleAttributes: eligible("B2-G"),
			Requirements:       sameForCourses("Required", allCourseValues...),
		}, unspecified},
		{"学部生でクラス未入力", testUser("B3", "InformationSystem", ""), noEligible, unspecified},
		{"1年生でクラス未入力", testUser("B1", "", ""), noEligible, unspecified},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requirementTypeFor(tt.user, tt.subject); got != tt.want {
				t.Errorf("requirementTypeFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
