package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type SyllabusHandler struct {
	c crud[model.Syllabus]
}

func NewSyllabusHandler(repo Repository[model.Syllabus]) *SyllabusHandler {
	return &SyllabusHandler{c: crud[model.Syllabus]{
		name: "syllabus",
		repo: repo,
		keys: func(m model.Syllabus) []keyField { return []keyField{{"id", m.ID}} },
	}}
}

type syllabusInput interface {
	GetId() string
	GetCredit() int32
	GetName() string
	GetEnName() string
	GetGrades() string
	GetFacultyNames() string
	GetPracticalHomeFacultyCategory() string
	GetMultiplePersonTeachingForm() string
	GetTeachingForm() string
	GetSummary() string
	GetLearningOutcomes() string
	GetAssignments() string
	GetEvaluationMethod() string
	GetTextbooks() string
	GetReferenceBooks() string
	GetPrerequisites() string
	GetPreLearning() string
	GetPostLearning() string
	GetNotes() string
	GetKeywords() string
	GetTargetCourses() string
	GetTargetAreas() string
	GetClassifications() string
	GetTeachingLanguage() string
	GetContentsAndSchedule() string
	GetTeachingAndExamForm() string
	GetDsopSubject() string
}

func syllabusFromProto(in syllabusInput, s fieldSet) (model.Syllabus, error) {
	var v violations
	var m model.Syllabus
	if s.has("id") {
		m.ID = requireString(&v, "id", in.GetId())
	}
	if s.has("name") {
		requireString(&v, "name", in.GetName())
	}
	if s.has("credit") {
		m.Credit = requireNonNegative(&v, "credit", in.GetCredit())
	}
	if s.has("name") {
		m.Name = in.GetName()
	}
	if s.has("en_name") {
		m.EnName = in.GetEnName()
	}
	if s.has("grades") {
		m.Grades = in.GetGrades()
	}
	if s.has("faculty_names") {
		m.FacultyNames = in.GetFacultyNames()
	}
	if s.has("practical_home_faculty_category") {
		m.PracticalHomeFacultyCategory = in.GetPracticalHomeFacultyCategory()
	}
	if s.has("multiple_person_teaching_form") {
		m.MultiplePersonTeachingForm = in.GetMultiplePersonTeachingForm()
	}
	if s.has("teaching_form") {
		m.TeachingForm = in.GetTeachingForm()
	}
	if s.has("summary") {
		m.Summary = in.GetSummary()
	}
	if s.has("learning_outcomes") {
		m.LearningOutcomes = in.GetLearningOutcomes()
	}
	if s.has("assignments") {
		m.Assignments = in.GetAssignments()
	}
	if s.has("evaluation_method") {
		m.EvaluationMethod = in.GetEvaluationMethod()
	}
	if s.has("textbooks") {
		m.Textbooks = in.GetTextbooks()
	}
	if s.has("reference_books") {
		m.ReferenceBooks = in.GetReferenceBooks()
	}
	if s.has("prerequisites") {
		m.Prerequisites = in.GetPrerequisites()
	}
	if s.has("pre_learning") {
		m.PreLearning = in.GetPreLearning()
	}
	if s.has("post_learning") {
		m.PostLearning = in.GetPostLearning()
	}
	if s.has("notes") {
		m.Notes = in.GetNotes()
	}
	if s.has("keywords") {
		m.Keywords = in.GetKeywords()
	}
	if s.has("target_courses") {
		m.TargetCourses = in.GetTargetCourses()
	}
	if s.has("target_areas") {
		m.TargetAreas = in.GetTargetAreas()
	}
	if s.has("classifications") {
		m.Classifications = in.GetClassifications()
	}
	if s.has("teaching_language") {
		m.TeachingLanguage = in.GetTeachingLanguage()
	}
	if s.has("contents_and_schedule") {
		m.ContentsAndSchedule = in.GetContentsAndSchedule()
	}
	if s.has("teaching_and_exam_form") {
		m.TeachingAndExamForm = in.GetTeachingAndExamForm()
	}
	if s.has("dsop_subject") {
		m.DsopSubject = in.GetDsopSubject()
	}
	return m, v.err()
}

func syllabusToProto(m model.Syllabus) *adminv1.Syllabus {
	return &adminv1.Syllabus{
		Id:                           m.ID,
		Credit:                       int32(m.Credit),
		Name:                         m.Name,
		EnName:                       m.EnName,
		Grades:                       m.Grades,
		FacultyNames:                 m.FacultyNames,
		PracticalHomeFacultyCategory: m.PracticalHomeFacultyCategory,
		MultiplePersonTeachingForm:   m.MultiplePersonTeachingForm,
		TeachingForm:                 m.TeachingForm,
		Summary:                      m.Summary,
		LearningOutcomes:             m.LearningOutcomes,
		Assignments:                  m.Assignments,
		EvaluationMethod:             m.EvaluationMethod,
		Textbooks:                    m.Textbooks,
		ReferenceBooks:               m.ReferenceBooks,
		Prerequisites:                m.Prerequisites,
		PreLearning:                  m.PreLearning,
		PostLearning:                 m.PostLearning,
		Notes:                        m.Notes,
		Keywords:                     m.Keywords,
		TargetCourses:                m.TargetCourses,
		TargetAreas:                  m.TargetAreas,
		Classifications:              m.Classifications,
		TeachingLanguage:             m.TeachingLanguage,
		ContentsAndSchedule:          m.ContentsAndSchedule,
		TeachingAndExamForm:          m.TeachingAndExamForm,
		DsopSubject:                  m.DsopSubject,
		CreatedAt:                    toTimestamp(m.CreatedAt),
		UpdatedAt:                    toTimestamp(m.UpdatedAt),
	}
}

func (h *SyllabusHandler) ListSyllabi(ctx context.Context, req *connect.Request[adminv1.ListSyllabiRequest]) (*connect.Response[adminv1.ListSyllabiResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListSyllabiResponse{NextPageToken: next}
	for _, m := range items {
		res.Syllabi = append(res.Syllabi, syllabusToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *SyllabusHandler) GetSyllabus(ctx context.Context, req *connect.Request[adminv1.GetSyllabusRequest]) (*connect.Response[adminv1.GetSyllabusResponse], error) {
	id, err := requireKeyString("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetSyllabusResponse{Syllabus: syllabusToProto(m)}), nil
}

func (h *SyllabusHandler) CreateSyllabus(ctx context.Context, req *connect.Request[adminv1.CreateSyllabusRequest]) (*connect.Response[adminv1.CreateSyllabusResponse], error) {
	m, err := syllabusFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateSyllabusResponse{Syllabus: syllabusToProto(created)}), nil
}

func (h *SyllabusHandler) UpdateSyllabus(ctx context.Context, req *connect.Request[adminv1.UpdateSyllabusRequest]) (*connect.Response[adminv1.UpdateSyllabusResponse], error) {
	id, err := requireKeyString("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := syllabusFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateSyllabusResponse{Syllabus: syllabusToProto(updated)}), nil
}

func (h *SyllabusHandler) DeleteSyllabus(ctx context.Context, req *connect.Request[adminv1.DeleteSyllabusRequest]) (*connect.Response[adminv1.DeleteSyllabusResponse], error) {
	id, err := requireKeyString("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteSyllabusResponse{}), nil
}
