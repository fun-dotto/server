package handler

import (
	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

// Subject

func subjectToAPI(d domain.Subject) api.Subject {
	eligible := make([]api.SubjectTargetClass, len(d.EligibleAttributes))
	for i, e := range d.EligibleAttributes {
		tc := api.SubjectTargetClass{
			Grade: api.DottoFoundationV1Grade(e.Grade),
		}
		if e.Class != nil {
			c := api.DottoFoundationV1Class(*e.Class)
			tc.Class = &c
		}
		eligible[i] = tc
	}

	requirements := make([]api.SubjectRequirement, len(d.Requirements))
	for i, r := range d.Requirements {
		requirements[i] = api.SubjectRequirement{
			Course:          api.DottoFoundationV1Course(r.Course.ID), // TODO: r.Course.ID はUUIDのためenumとして不正。domain.Course にenum値のフィールドを追加して正しい値を使う
			RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
		}
	}

	faculties := make([]api.SubjectFaculty, 0)
	// TODO: domain model から faculties を取得して設定する

	return api.Subject{
		Id:                 d.ID,
		Name:               d.Name,
		Credit:             0, // TODO: domain.Subject に Credit フィールドを追加してから取得する
		Faculties:          faculties,
		Semester:           api.DottoFoundationV1CourseSemester(d.Semester),
		EligibleAttributes: eligible,
		Requirements:       requirements,
	}
}

func subjectToSummaryAPI(d domain.Subject) api.SubjectSummary {
	faculties := make([]api.SubjectFaculty, 0)
	// TODO: domain model から faculties を取得して設定する

	return api.SubjectSummary{
		Id:        d.ID,
		Name:      d.Name,
		Faculties: faculties,
	}
}

func subjectsToSummaryAPI(ds []domain.Subject) []api.SubjectSummary {
	result := make([]api.SubjectSummary, len(ds))
	for i, d := range ds {
		result[i] = subjectToSummaryAPI(d)
	}
	return result
}
