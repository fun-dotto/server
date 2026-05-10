package handler

import (
	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func syllabusToAPI(d domain.Syllabus) api.Syllabus {
	return api.Syllabus{
		Id:                           d.ID,
		Name:                         d.Name,
		EnName:                       d.EnName,
		Grades:                       d.Grades,
		Credit:                       d.Credit,
		FacultyNames:                 d.FacultyNames,
		PracticalHomeFacultyCategory: d.PracticalHomeFacultyCategory,
		MultiplePersonTeachingForm:   d.MultiplePersonTeachingForm,
		TeachingForm:                 d.TeachingForm,
		Summary:                      d.Summary,
		LearningOutcomes:             d.LearningOutcomes,
		Assignments:                  d.Assignments,
		EvaluationMethod:             d.EvaluationMethod,
		Textbooks:                    d.Textbooks,
		ReferenceBooks:               d.ReferenceBooks,
		Prerequisites:                d.Prerequisites,
		PreLearning:                  d.PreLearning,
		PostLearning:                 d.PostLearning,
		Notes:                        d.Notes,
		Keywords:                     d.Keywords,
		TargetCourses:                d.TargetCourses,
		TargetAreas:                  d.TargetAreas,
		Classifications:              d.Classifications,
		TeachingLanguage:             d.TeachingLanguage,
		ContentsAndSchedule:          d.ContentsAndSchedule,
		TeachingAndExamForm:          d.TeachingAndExamForm,
		DsopSubject:                  d.DsopSubject,
	}
}
