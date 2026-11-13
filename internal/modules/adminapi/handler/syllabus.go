package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SyllabusHandler struct {
	c crud[model.Syllabus]
}

func NewSyllabusHandler(db *gorm.DB) *SyllabusHandler {
	return &SyllabusHandler{c: crud[model.Syllabus]{
		name:     "syllabus",
		db:       db,
		keys:     func(m model.Syllabus) []keyField { return []keyField{{"id", m.ID.String()}} },
		updateFn: updateSyllabus,
	}}
}

// updateSyllabus は columns を更新したあと、更新後の内容から content_hash を計算し直す。
func updateSyllabus(db *gorm.DB, key map[string]any, v *model.Syllabus, columns []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := updateColumns(tx, key, v, columns); err != nil {
			return err
		}
		var current model.Syllabus
		if err := tx.Where(key).Take(&current).Error; err != nil {
			return err
		}
		return tx.Model(&model.Syllabus{}).Where(key).Omit(clause.Associations).
			Update("content_hash", current.ComputeContentHash()).Error
	})
}

type syllabusInput interface {
	GetLessonId() string
	GetYear() int32
	GetDisplayOrder() string
	GetListLessonName() string
	GetListEnName() string
	GetListStaffName() string
	GetListTeachingForm() string
	GetName() string
	GetEnName() string
	GetGrades() string
	GetTerm() string
	GetCredit() string
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
	GetClassifications() string
	GetDsopSubject() string
	GetTeachingLanguage() string
	GetTeachingAndExamForm() string
	GetContentsAndSchedule() string
	GetEnSummary() string
	GetEnLearningOutcomes() string
	GetEnContentsAndSchedule() string
	GetEnEvaluationMethod() string
	GetEnTeachingLanguage() string
	GetEnPrerequisites() string
	GetEnPrePostLearning() string
	GetEnTextbooks() string
	GetEnNotes() string
	GetEnTeachingAndExamForm() string
	GetDiplomaPolicies() string
	GetDiplomaPolicyCorrespondence() string
}

func syllabusFromProto(in syllabusInput, s fieldSet) (model.Syllabus, error) {
	var v violations
	var m model.Syllabus
	if s.has("lesson_id") {
		m.LessonID = requireString(&v, "lesson_id", in.GetLessonId())
	}
	if s.has("year") {
		m.Year = requirePositive(&v, "year", in.GetYear())
	}
	if s.has("name") {
		requireString(&v, "name", in.GetName())
	}
	if s.has("display_order") {
		m.DisplayOrder = in.GetDisplayOrder()
	}
	if s.has("list_lesson_name") {
		m.ListLessonName = in.GetListLessonName()
	}
	if s.has("list_en_name") {
		m.ListEnName = in.GetListEnName()
	}
	if s.has("list_staff_name") {
		m.ListStaffName = in.GetListStaffName()
	}
	if s.has("list_teaching_form") {
		m.ListTeachingForm = in.GetListTeachingForm()
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
	if s.has("term") {
		m.Term = in.GetTerm()
	}
	if s.has("credit") {
		m.Credit = in.GetCredit()
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
	if s.has("classifications") {
		m.Classifications = in.GetClassifications()
	}
	if s.has("dsop_subject") {
		m.DsopSubject = in.GetDsopSubject()
	}
	if s.has("teaching_language") {
		m.TeachingLanguage = in.GetTeachingLanguage()
	}
	if s.has("teaching_and_exam_form") {
		m.TeachingAndExamForm = in.GetTeachingAndExamForm()
	}
	if s.has("contents_and_schedule") {
		m.ContentsAndSchedule = in.GetContentsAndSchedule()
	}
	if s.has("en_summary") {
		m.EnSummary = in.GetEnSummary()
	}
	if s.has("en_learning_outcomes") {
		m.EnLearningOutcomes = in.GetEnLearningOutcomes()
	}
	if s.has("en_contents_and_schedule") {
		m.EnContentsAndSchedule = in.GetEnContentsAndSchedule()
	}
	if s.has("en_evaluation_method") {
		m.EnEvaluationMethod = in.GetEnEvaluationMethod()
	}
	if s.has("en_teaching_language") {
		m.EnTeachingLanguage = in.GetEnTeachingLanguage()
	}
	if s.has("en_prerequisites") {
		m.EnPrerequisites = in.GetEnPrerequisites()
	}
	if s.has("en_pre_post_learning") {
		m.EnPrePostLearning = in.GetEnPrePostLearning()
	}
	if s.has("en_textbooks") {
		m.EnTextbooks = in.GetEnTextbooks()
	}
	if s.has("en_notes") {
		m.EnNotes = in.GetEnNotes()
	}
	if s.has("en_teaching_and_exam_form") {
		m.EnTeachingAndExamForm = in.GetEnTeachingAndExamForm()
	}
	if s.has("diploma_policies") {
		m.DiplomaPolicies = in.GetDiplomaPolicies()
	}
	if s.has("diploma_policy_correspondence") {
		m.DiplomaPolicyCorrespondence = in.GetDiplomaPolicyCorrespondence()
	}
	return m, v.err()
}

func syllabusToProto(m model.Syllabus) *adminv1.Syllabus {
	return &adminv1.Syllabus{
		Id:                           m.ID.String(),
		LessonId:                     m.LessonID,
		Year:                         int32(m.Year),
		DisplayOrder:                 m.DisplayOrder,
		ListLessonName:               m.ListLessonName,
		ListEnName:                   m.ListEnName,
		ListStaffName:                m.ListStaffName,
		ListTeachingForm:             m.ListTeachingForm,
		Name:                         m.Name,
		EnName:                       m.EnName,
		Grades:                       m.Grades,
		Term:                         m.Term,
		Credit:                       m.Credit,
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
		Classifications:              m.Classifications,
		DsopSubject:                  m.DsopSubject,
		TeachingLanguage:             m.TeachingLanguage,
		TeachingAndExamForm:          m.TeachingAndExamForm,
		ContentsAndSchedule:          m.ContentsAndSchedule,
		EnSummary:                    m.EnSummary,
		EnLearningOutcomes:           m.EnLearningOutcomes,
		EnContentsAndSchedule:        m.EnContentsAndSchedule,
		EnEvaluationMethod:           m.EnEvaluationMethod,
		EnTeachingLanguage:           m.EnTeachingLanguage,
		EnPrerequisites:              m.EnPrerequisites,
		EnPrePostLearning:            m.EnPrePostLearning,
		EnTextbooks:                  m.EnTextbooks,
		EnNotes:                      m.EnNotes,
		EnTeachingAndExamForm:        m.EnTeachingAndExamForm,
		DiplomaPolicies:              m.DiplomaPolicies,
		DiplomaPolicyCorrespondence:  m.DiplomaPolicyCorrespondence,
		ContentHash:                  m.ContentHash,
		FetchedAt:                    toTimestamp(m.FetchedAt),
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
	id, err := parseIDField("id", req.Msg.GetId())
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
	m.ContentHash = m.ComputeContentHash()
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateSyllabusResponse{Syllabus: syllabusToProto(created)}), nil
}

func (h *SyllabusHandler) UpdateSyllabus(ctx context.Context, req *connect.Request[adminv1.UpdateSyllabusRequest]) (*connect.Response[adminv1.UpdateSyllabusResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
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
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteSyllabusResponse{}), nil
}
