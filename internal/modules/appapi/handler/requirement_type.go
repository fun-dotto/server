package handler

import (
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

const unspecifiedRequirementType = appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_UNSPECIFIED

// gradeB1 は学部 1 年の学年の DB 値。
const gradeB1 = "B1"

// undergraduateGrades は学部生の学年の DB 値。
var undergraduateGrades = map[string]bool{"B1": true, "B2": true, "B3": true, "B4": true}

// requirementTypeFor は、ユーザーにとっての科目の必修区分を返す。
// 欠けた情報は推測せず、判定できない場合は UNSPECIFIED を返す。上から順に判定する。
//
//  1. ユーザーが存在しない（user が nil）: UNSPECIFIED
//  2. 登録情報に未入力がある（hasMissingProfile）: UNSPECIFIED
//  3. 科目の対象学年・クラスにユーザーが含まれない（isEligible）: UNSPECIFIED
//  4. コースが設定されている: そのコースの区分
//  5. コースが未設定（ここに来るのは 1 年生だけ）: 全コースで同じ区分ならその区分
//
// 再履修（対象学年でない科目の履修）は考慮しない。
func requirementTypeFor(user *model.User, subject *model.Subject) appv1.SubjectRequirementType {
	if user == nil || subject == nil || hasMissingProfile(user) {
		return unspecifiedRequirementType
	}
	if !isEligible(user, subject.EligibleAttributes) {
		return unspecifiedRequirementType
	}
	if user.Course != nil {
		return requirementTypeOfCourse(subject.Requirements, *user.Course)
	}
	return commonRequirementType(subject.Requirements)
}

// hasMissingProfile は、ユーザーの登録情報（学年・コース・クラス）に未入力の項目があるかを返す。
// 「入力していない」と「なしを選んだ」は区別できないので、どちらも未入力とする。ただし次は未入力にしない。
//   - 1 年生のコース: コース配属は 2 年次なので、未設定が通常の状態
//   - 院生のクラス: 院生向けの科目は、学年だけで対象を決めている
func hasMissingProfile(u *model.User) bool {
	if u.Grade == nil {
		return true
	}
	if u.Course == nil && *u.Grade != gradeB1 {
		return true
	}
	return u.Class == nil && undergraduateGrades[*u.Grade]
}

// isEligible は、ユーザーが科目の対象学年・クラスに含まれるかを返す。
// 対象学年・クラスが空の科目は、全員を対象とする。クラスが未設定の行は、学年だけで判定する。
func isEligible(u *model.User, attrs []model.SubjectEligibleAttribute) bool {
	if len(attrs) == 0 {
		return true
	}
	if u.Grade == nil {
		return false
	}
	for _, a := range attrs {
		if a.Grade != *u.Grade {
			continue
		}
		if a.Class == nil || (u.Class != nil && *a.Class == *u.Class) {
			return true
		}
	}
	return false
}

// requirementTypeOfCourse は、科目のそのコースの区分を返す。
// 行がない場合と、同じコースに異なる区分の行がある場合は UNSPECIFIED。
func requirementTypeOfCourse(reqs []model.SubjectRequirement, course string) appv1.SubjectRequirementType {
	result := unspecifiedRequirementType
	found := false
	for _, r := range reqs {
		if r.Course != course {
			continue
		}
		t := requirementTypeMap.toProto[r.RequirementType]
		if found && t != result {
			return unspecifiedRequirementType
		}
		result, found = t, true
	}
	return result
}

// commonRequirementType は、全コースで区分が同じ場合にその区分を返す。
// 行のないコースがある場合と、コースによって区分が異なる場合は UNSPECIFIED。
func commonRequirementType(reqs []model.SubjectRequirement) appv1.SubjectRequirementType {
	result := unspecifiedRequirementType
	for course := range courseMap.toProto {
		t := requirementTypeOfCourse(reqs, course)
		if t == unspecifiedRequirementType {
			return unspecifiedRequirementType
		}
		if result != unspecifiedRequirementType && t != result {
			return unspecifiedRequirementType
		}
		result = t
	}
	return result
}
