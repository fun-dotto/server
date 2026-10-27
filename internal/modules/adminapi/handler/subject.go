package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type SubjectHandler struct {
	c crud[model.Subject]
}

func NewSubjectHandler(repo Repository[model.Subject]) *SubjectHandler {
	return &SubjectHandler{c: crud[model.Subject]{
		name: "subject",
		repo: repo,
		keys: func(m model.Subject) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

var subjectSemesterEnum = enumMap[adminv1.SubjectSemester]{
	adminv1.SubjectSemester_SUBJECT_SEMESTER_ALL_YEAR:         "AllYear",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_H1:               "H1",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_H2:               "H2",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_Q1:               "Q1",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_Q2:               "Q2",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_Q3:               "Q3",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_Q4:               "Q4",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_SUMMER_INTENSIVE: "SummerIntensive",
	adminv1.SubjectSemester_SUBJECT_SEMESTER_WINTER_INTENSIVE: "WinterIntensive",
}

var subjectClassificationEnum = enumMap[adminv1.SubjectClassification]{
	adminv1.SubjectClassification_SUBJECT_CLASSIFICATION_SPECIALIZED:          "Specialized",
	adminv1.SubjectClassification_SUBJECT_CLASSIFICATION_CULTURAL:             "Cultural",
	adminv1.SubjectClassification_SUBJECT_CLASSIFICATION_RESEARCH_INSTRUCTION: "ResearchInstruction",
}

var subjectCulturalSubjectCategoryEnum = enumMap[adminv1.SubjectCulturalSubjectCategory]{
	adminv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_SOCIETY:       "Society",
	adminv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_HUMAN:         "Human",
	adminv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_SCIENCE:       "Science",
	adminv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_HEALTH:        "Health",
	adminv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_COMMUNICATION: "Communication",
}

var subjectRequirementTypeEnum = enumMap[adminv1.SubjectRequirementType]{
	adminv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_REQUIRED:          "Required",
	adminv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL:          "Optional",
	adminv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL_REQUIRED: "OptionalRequired",
}

// subjectColumns は子テーブルとして保存するフィールド。repository.SubjectRepository が置き換える。
var subjectColumns = map[string][]string{
	"faculties":           {"Faculties"},
	"eligible_attributes": {"EligibleAttributes"},
	"requirements":        {"Requirements"},
}

type subjectInput interface {
	GetName() string
	GetYear() int32
	GetSemester() adminv1.SubjectSemester
	GetCredit() int32
	GetClassification() adminv1.SubjectClassification
	GetCulturalSubjectCategory() adminv1.SubjectCulturalSubjectCategory
	GetSyllabusId() string
	GetFaculties() []*adminv1.SubjectFaculty
	GetEligibleAttributes() []*adminv1.SubjectEligibleAttribute
	GetRequirements() []*adminv1.SubjectRequirement
}

func subjectFromProto(in subjectInput, s fieldSet) (model.Subject, error) {
	var v violations
	var m model.Subject
	if s.has("name") {
		m.Name = requireString(&v, "name", in.GetName())
	}
	if s.has("year") {
		m.Year = requirePositive(&v, "year", in.GetYear())
	}
	if s.has("semester") {
		m.Semester = subjectSemesterEnum.toDB(&v, "semester", in.GetSemester(), true)
	}
	if s.has("credit") {
		m.Credit = requireNonNegative(&v, "credit", in.GetCredit())
	}
	if s.has("classification") {
		m.Classification = subjectClassificationEnum.toDB(&v, "classification", in.GetClassification(), true)
	}
	if s.has("cultural_subject_category") {
		// 教養科目以外では未指定（空文字）を許す
		m.CulturalSubjectCategory = subjectCulturalSubjectCategoryEnum.toDB(&v, "cultural_subject_category", in.GetCulturalSubjectCategory(), false)
	}
	if s.has("syllabus_id") {
		m.SyllabusID = requireString(&v, "syllabus_id", in.GetSyllabusId())
	}
	if s.has("faculties") {
		for i, f := range in.GetFaculties() {
			m.Faculties = append(m.Faculties, model.SubjectFaculty{
				FacultyID: parseUUID(&v, fmt.Sprintf("faculties.%d.faculty_id", i), f.GetFacultyId()),
				IsPrimary: f.GetIsPrimary(),
			})
		}
	}
	if s.has("eligible_attributes") {
		for i, a := range in.GetEligibleAttributes() {
			attr := model.SubjectEligibleAttribute{
				Grade: gradeEnum.toDB(&v, fmt.Sprintf("eligible_attributes.%d.grade", i), a.GetGrade(), true),
			}
			if a.Class != nil {
				attr.Class = classEnum.toOptionalDB(&v, fmt.Sprintf("eligible_attributes.%d.class", i), a.Class)
			}
			m.EligibleAttributes = append(m.EligibleAttributes, attr)
		}
	}
	if s.has("requirements") {
		for i, r := range in.GetRequirements() {
			m.Requirements = append(m.Requirements, model.SubjectRequirement{
				Course:          courseEnum.toDB(&v, fmt.Sprintf("requirements.%d.course", i), r.GetCourse(), true),
				RequirementType: subjectRequirementTypeEnum.toDB(&v, fmt.Sprintf("requirements.%d.requirement_type", i), r.GetRequirementType(), true),
			})
		}
	}
	return m, v.err()
}

func subjectToProto(m model.Subject) *adminv1.Subject {
	s := &adminv1.Subject{
		Id:                      m.ID.String(),
		Name:                    m.Name,
		Year:                    int32(m.Year),
		Semester:                subjectSemesterEnum.fromDB(m.Semester),
		Credit:                  int32(m.Credit),
		Classification:          subjectClassificationEnum.fromDB(m.Classification),
		CulturalSubjectCategory: subjectCulturalSubjectCategoryEnum.fromDB(m.CulturalSubjectCategory),
		SyllabusId:              m.SyllabusID,
		CreatedAt:               toTimestamp(m.CreatedAt),
		UpdatedAt:               toTimestamp(m.UpdatedAt),
	}
	for _, f := range m.Faculties {
		s.Faculties = append(s.Faculties, &adminv1.SubjectFaculty{FacultyId: f.FacultyID.String(), IsPrimary: f.IsPrimary})
	}
	for _, a := range m.EligibleAttributes {
		s.EligibleAttributes = append(s.EligibleAttributes, &adminv1.SubjectEligibleAttribute{
			Grade: gradeEnum.fromDB(a.Grade),
			Class: classEnum.fromOptionalDB(a.Class),
		})
	}
	for _, r := range m.Requirements {
		s.Requirements = append(s.Requirements, &adminv1.SubjectRequirement{
			Course:          courseEnum.fromDB(r.Course),
			RequirementType: subjectRequirementTypeEnum.fromDB(r.RequirementType),
		})
	}
	return s
}

func (h *SubjectHandler) ListSubjects(ctx context.Context, req *connect.Request[adminv1.ListSubjectsRequest]) (*connect.Response[adminv1.ListSubjectsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListSubjectsResponse{NextPageToken: next}
	for _, m := range items {
		res.Subjects = append(res.Subjects, subjectToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *SubjectHandler) GetSubject(ctx context.Context, req *connect.Request[adminv1.GetSubjectRequest]) (*connect.Response[adminv1.GetSubjectResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetSubjectResponse{Subject: subjectToProto(m)}), nil
}

func (h *SubjectHandler) CreateSubject(ctx context.Context, req *connect.Request[adminv1.CreateSubjectRequest]) (*connect.Response[adminv1.CreateSubjectResponse], error) {
	m, err := subjectFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateSubjectResponse{Subject: subjectToProto(created)}), nil
}

func (h *SubjectHandler) UpdateSubject(ctx context.Context, req *connect.Request[adminv1.UpdateSubjectRequest]) (*connect.Response[adminv1.UpdateSubjectResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := subjectFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(subjectColumns))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateSubjectResponse{Subject: subjectToProto(updated)}), nil
}

func (h *SubjectHandler) DeleteSubject(ctx context.Context, req *connect.Request[adminv1.DeleteSubjectRequest]) (*connect.Response[adminv1.DeleteSubjectResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteSubjectResponse{}), nil
}
