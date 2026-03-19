package handler

import (
	"time"

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
			Course:          api.DottoFoundationV1Course(r.Course),
			RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
		}
	}

	faculties := make([]api.SubjectFaculty, len(d.Faculties))
	for i, f := range d.Faculties {
		// TODO: Faculty には Name, Email も必要。現状 domain.SubjectFaculty に情報がないため Id のみセットしている。
		faculties[i] = api.SubjectFaculty{
			Faculty:   api.Faculty{Id: f.FacultyID},
			IsPrimary: f.IsPrimary,
		}
	}

	return api.Subject{
		Id:                 d.ID,
		Name:               d.Name,
		Faculties:          faculties,
		Year:               d.Year,
		Semester:           api.DottoFoundationV1CourseSemester(d.Semester),
		Credit:             d.Credit,
		EligibleAttributes: eligible,
		Requirements:       requirements,
	}
}

func subjectToSummaryAPI(d domain.Subject) api.SubjectSummary {
	faculties := make([]api.SubjectFaculty, len(d.Faculties))
	for i, f := range d.Faculties {
		// TODO: Faculty には Name, Email も必要。現状 domain.SubjectFaculty に情報がないため Id のみセットしている。
		faculties[i] = api.SubjectFaculty{
			Faculty:   api.Faculty{Id: f.FacultyID},
			IsPrimary: f.IsPrimary,
		}
	}

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

func buildSubjectListFilter(params api.SubjectsV1ListParams) domain.SubjectListFilter {
	filter := domain.SubjectListFilter{}

	if params.Ids != nil {
		filter.IDs = *params.Ids
	}
	if params.Q != nil {
		filter.Q = params.Q
	}
	if params.Grade != nil {
		grades := make([]domain.Grade, len(*params.Grade))
		for i, g := range *params.Grade {
			grades[i] = domain.Grade(g)
		}
		filter.Grade = grades
	}
	if params.Courses != nil {
		courses := make([]domain.CourseType, len(*params.Courses))
		for i, c := range *params.Courses {
			courses[i] = domain.CourseType(c)
		}
		filter.Courses = courses
	}
	if params.Class != nil {
		classes := make([]domain.Class, len(*params.Class))
		for i, c := range *params.Class {
			classes[i] = domain.Class(c)
		}
		filter.Class = classes
	}
	if params.Classification != nil {
		classifications := make([]domain.SubjectClassification, len(*params.Classification))
		for i, c := range *params.Classification {
			classifications[i] = domain.SubjectClassification(c)
		}
		filter.Classification = classifications
	}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		// デフォルトで今年度を設定
		// TODO: このデフォルト値設定のロジックは service 層に移すべき。
		// また、日本の大学の年度は4月始まりのため、1〜3月は前年度を返す必要がある。
		currentYear := time.Now().Year()
		filter.Year = &currentYear
	}
	if params.Semester != nil {
		semesters := make([]domain.CourseSemester, len(*params.Semester))
		for i, s := range *params.Semester {
			semesters[i] = domain.CourseSemester(s)
		}
		filter.Semester = semesters
	}
	if params.RequirementType != nil {
		reqTypes := make([]domain.SubjectRequirementType, len(*params.RequirementType))
		for i, r := range *params.RequirementType {
			reqTypes[i] = domain.SubjectRequirementType(r)
		}
		filter.RequirementType = reqTypes
	}
	if params.CulturalSubjectCategory != nil {
		cats := make([]domain.CulturalSubjectCategory, len(*params.CulturalSubjectCategory))
		for i, c := range *params.CulturalSubjectCategory {
			cats[i] = domain.CulturalSubjectCategory(c)
		}
		filter.CulturalSubjectCategory = cats
	}

	return filter
}

// Syllabus

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

func facultyToAPI(faculty domain.Faculty) api.Faculty {
	return api.Faculty{
		Id:    faculty.ID,
		Name:  faculty.Name,
		Email: faculty.Email,
	}
}

func facultiesToAPI(faculties []domain.Faculty) []api.Faculty {
	result := make([]api.Faculty, len(faculties))
	for i, faculty := range faculties {
		result[i] = facultyToAPI(faculty)
	}
	return result
}

func toDomainFacultyFromRequest(id string, req api.FacultyRequest) domain.Faculty {
	return domain.Faculty{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
}
