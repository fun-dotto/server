package handler

import (
	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func syllabusToAPI(d domain.Syllabus) api.Syllabus {
	return api.Syllabus{
		Id:                           d.ID,
		LessonId:                     d.LessonID,
		Year:                         d.Year,
		DisplayOrder:                 d.DisplayOrder,
		ListLessonName:               d.ListLessonName,
		ListEnName:                   d.ListEnName,
		ListStaffName:                d.ListStaffName,
		ListTeachingForm:             d.ListTeachingForm,
		Name:                         d.Name,
		EnName:                       d.EnName,
		Grades:                       d.Grades,
		Term:                         d.Term,
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
		Classifications:              d.Classifications,
		DsopSubject:                  d.DsopSubject,
		TeachingLanguage:             d.TeachingLanguage,
		TeachingAndExamForm:          d.TeachingAndExamForm,
		ContentsAndSchedule:          d.ContentsAndSchedule,
		EnSummary:                    d.EnSummary,
		EnLearningOutcomes:           d.EnLearningOutcomes,
		EnContentsAndSchedule:        d.EnContentsAndSchedule,
		EnEvaluationMethod:           d.EnEvaluationMethod,
		EnTeachingLanguage:           d.EnTeachingLanguage,
		EnPrerequisites:              d.EnPrerequisites,
		EnPrePostLearning:            d.EnPrePostLearning,
		EnTextbooks:                  d.EnTextbooks,
		EnNotes:                      d.EnNotes,
		EnTeachingAndExamForm:        d.EnTeachingAndExamForm,
		DiplomaPolicies:              d.DiplomaPolicies,
		DiplomaPolicyCorrespondence:  d.DiplomaPolicyCorrespondence,
		ContentHash:                  d.ContentHash,
		FetchedAt:                    d.FetchedAt,
	}
}
