package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

func (h *Handler) SubjectsV1List(ctx context.Context, request api.SubjectsV1ListRequestObject) (api.SubjectsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toSubjectQuery(request.Params)

	subjects, err := h.academicService.GetSubjects(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects: %w", err)
	}

	apiSubjects := make([]api.SubjectSummary, len(subjects))
	for i, subject := range subjects {
		apiSubjects[i] = toApiSubjectSummary(subject)
	}

	return api.SubjectsV1List200JSONResponse{
		Subjects: apiSubjects,
	}, nil
}

func (h *Handler) SubjectsV1Detail(ctx context.Context, request api.SubjectsV1DetailRequestObject) (api.SubjectsV1DetailResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	subject, err := h.academicService.GetSubject(request.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject: %w", err)
	}

	return api.SubjectsV1Detail200JSONResponse{
		Subject: toApiSubjectDetail(*subject),
	}, nil
}

// toSubjectQuery は SubjectsV1List の API パラメータを domain.SubjectQuery に変換する
func toSubjectQuery(params api.SubjectsV1ListParams) domain.SubjectQuery {
	var q string
	if params.Q != nil {
		q = *params.Q
	}

	var grade []api.DottoFoundationV1Grade
	if params.Grades != nil {
		grade = *params.Grades
	}

	var courses []api.DottoFoundationV1Course
	if params.Courses != nil {
		courses = *params.Courses
	}

	var class []api.DottoFoundationV1Class
	if params.Classes != nil {
		class = *params.Classes
	}

	var classification []api.DottoFoundationV1SubjectClassification
	if params.Classifications != nil {
		classification = *params.Classifications
	}

	var semester []api.DottoFoundationV1CourseSemester
	if params.Semesters != nil {
		semester = *params.Semesters
	}

	var requirementType []api.DottoFoundationV1SubjectRequirementType
	if params.RequirementTypes != nil {
		requirementType = *params.RequirementTypes
	}

	var culturalSubjectCategory []api.DottoFoundationV1CulturalSubjectCategory
	if params.CulturalSubjectCategories != nil {
		culturalSubjectCategory = *params.CulturalSubjectCategories
	}

	return domain.SubjectQuery{
		Q:                       q,
		Grade:                   toGrades(grade),
		Courses:                 toCourses(courses),
		Class:                   toClasses(class),
		Classification:          toClassifications(classification),
		Semester:                toSemesters(semester),
		RequirementType:         toRequirementTypes(requirementType),
		CulturalSubjectCategory: toCulturalSubjectCategories(culturalSubjectCategory),
	}
}

// toGrades はAPIの学年をDomainの学年に変換する
func toGrades(grades []api.DottoFoundationV1Grade) []domain.Grade {
	result := make([]domain.Grade, len(grades))
	for i, g := range grades {
		result[i] = domain.Grade(g)
	}
	return result
}

// toCourses はAPIのコースをDomainのコースに変換する
func toCourses(courses []api.DottoFoundationV1Course) []domain.Course {
	result := make([]domain.Course, len(courses))
	for i, c := range courses {
		result[i] = domain.Course(c)
	}
	return result
}

// toClasses はAPIのクラスをDomainのクラスに変換する
func toClasses(classes []api.DottoFoundationV1Class) []domain.Class {
	result := make([]domain.Class, len(classes))
	for i, c := range classes {
		result[i] = domain.Class(c)
	}
	return result
}

// toClassifications はAPIの科目カテゴリをDomainの科目カテゴリに変換する
func toClassifications(classifications []api.DottoFoundationV1SubjectClassification) []domain.SubjectClassification {
	result := make([]domain.SubjectClassification, len(classifications))
	for i, c := range classifications {
		result[i] = domain.SubjectClassification(c)
	}
	return result
}

// toSemesters はAPIの開講時期をDomainの開講時期に変換する
func toSemesters(semesters []api.DottoFoundationV1CourseSemester) []domain.CourseSemester {
	result := make([]domain.CourseSemester, len(semesters))
	for i, s := range semesters {
		result[i] = domain.CourseSemester(s)
	}
	return result
}

// toRequirementTypes はAPIの必修・選択をDomainの必修・選択に変換する
func toRequirementTypes(types []api.DottoFoundationV1SubjectRequirementType) []domain.SubjectRequirementType {
	result := make([]domain.SubjectRequirementType, len(types))
	for i, t := range types {
		result[i] = domain.SubjectRequirementType(t)
	}
	return result
}

// toCulturalSubjectCategories はAPIの教養科目カテゴリをDomainの教養科目カテゴリに変換する
func toCulturalSubjectCategories(categories []api.DottoFoundationV1CulturalSubjectCategory) []domain.CulturalSubjectCategory {
	result := make([]domain.CulturalSubjectCategory, len(categories))
	for i, c := range categories {
		result[i] = domain.CulturalSubjectCategory(c)
	}
	return result
}

// toApiSubjectSummary はDomainの科目をAPIの科目サマリーに変換する
func toApiSubjectSummary(subject domain.Subject) api.SubjectSummary {
	return api.SubjectSummary{
		Id:        subject.ID,
		Name:      subject.Name,
		Year:      subject.Year,
		Semester:  api.DottoFoundationV1CourseSemester(subject.Semester),
		Credit:    subject.Credit,
		Faculties: toApiFaculties(subject.Faculties),
	}
}

// toApiSubjectDetail はDomainの科目をAPIの科目に変換する
func toApiSubjectDetail(subject domain.Subject) api.SubjectDetail {
	var syllabus api.AcademicServiceSyllabus
	if subject.Syllabus != nil {
		syllabus = toApiSyllabus(*subject.Syllabus)
	}

	return api.SubjectDetail{
		Id:                 subject.ID,
		Name:               subject.Name,
		Year:               subject.Year,
		Credit:             subject.Credit,
		Semester:           api.DottoFoundationV1CourseSemester(subject.Semester),
		Faculties:          toApiFaculties(subject.Faculties),
		Requirements:       toApiRequirements(subject.Requirements),
		EligibleAttributes: toApiTargetClasses(subject.EligibleAttributes),
		Syllabus:           syllabus,
	}
}

// toApiSyllabus はDomainのシラバスをAPIのシラバスに変換する
func toApiSyllabus(syllabus domain.Syllabus) api.AcademicServiceSyllabus {
	return api.AcademicServiceSyllabus{
		Id:                           syllabus.ID,
		Name:                         syllabus.Name,
		EnName:                       syllabus.EnName,
		Summary:                      syllabus.Summary,
		LearningOutcomes:             syllabus.LearningOutcomes,
		ContentsAndSchedule:          syllabus.ContentsAndSchedule,
		PreLearning:                  syllabus.PreLearning,
		PostLearning:                 syllabus.PostLearning,
		Assignments:                  syllabus.Assignments,
		EvaluationMethod:             syllabus.EvaluationMethod,
		Textbooks:                    syllabus.Textbooks,
		ReferenceBooks:               syllabus.ReferenceBooks,
		Prerequisites:                syllabus.Prerequisites,
		Notes:                        syllabus.Notes,
		Keywords:                     syllabus.Keywords,
		Classifications:              syllabus.Classifications,
		Grades:                       syllabus.Grades,
		Credit:                       syllabus.Credit,
		FacultyNames:                 syllabus.FacultyNames,
		TeachingForm:                 syllabus.TeachingForm,
		TeachingAndExamForm:          syllabus.TeachingAndExamForm,
		TeachingLanguage:             syllabus.TeachingLanguage,
		MultiplePersonTeachingForm:   syllabus.MultiplePersonTeachingForm,
		PracticalHomeFacultyCategory: syllabus.PracticalHomeFacultyCategory,
		DsopSubject:                  syllabus.DsopSubject,
		TargetAreas:                  syllabus.TargetAreas,
		TargetCourses:                syllabus.TargetCourses,
	}
}

// toApiFaculties はDomainの教員をAPIの教員に変換する
func toApiFaculties(faculties []domain.SubjectFaculty) []api.SubjectFaculty {
	result := make([]api.SubjectFaculty, len(faculties))
	for i, f := range faculties {
		result[i] = api.SubjectFaculty{
			Faculty: api.AcademicServiceFaculty{
				Id:    f.Faculty.ID,
				Name:  f.Faculty.Name,
				Email: f.Faculty.Email,
			},
			IsPrimary: f.IsPrimary,
		}
	}
	return result
}

// toApiRequirements はDomainの科目群・科目区分をAPIの科目群・科目区分に変換する
func toApiRequirements(requirements []domain.SubjectRequirement) []api.AcademicServiceSubjectRequirement {
	result := make([]api.AcademicServiceSubjectRequirement, len(requirements))
	for i, r := range requirements {
		result[i] = api.AcademicServiceSubjectRequirement{
			Course:          api.DottoFoundationV1Course(r.Course),
			RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
		}
	}
	return result
}

// toApiTargetClasses はDomainの対象学年・クラスをAPIの対象学年・クラスに変換する
func toApiTargetClasses(targetClasses []domain.SubjectTargetClass) []api.AcademicServiceSubjectTargetClass {
	result := make([]api.AcademicServiceSubjectTargetClass, len(targetClasses))
	for i, tc := range targetClasses {
		var class *api.DottoFoundationV1Class
		if tc.Class != nil {
			c := api.DottoFoundationV1Class(*tc.Class)
			class = &c
		}
		result[i] = api.AcademicServiceSubjectTargetClass{
			Grade: api.DottoFoundationV1Grade(tc.Grade),
			Class: class,
		}
	}
	return result
}
