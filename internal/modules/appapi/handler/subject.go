package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// SubjectHandler は SubjectService を実装する。
type SubjectHandler struct {
	appv1connect.UnimplementedSubjectServiceHandler
	service AcademicService
}

// NewSubjectHandler は SubjectHandler を返す。
func NewSubjectHandler(s AcademicService) *SubjectHandler {
	return &SubjectHandler{service: s}
}

// ListSubjects は科目を検索する。
func (h *SubjectHandler) ListSubjects(ctx context.Context, req *connect.Request[appv1.ListSubjectsRequest]) (*connect.Response[appv1.ListSubjectsResponse], error) {
	m := req.Msg
	var v violations
	query := domain.SubjectQuery{
		Q:                       m.GetQuery(),
		Grade:                   fromProtoEnums(gradeMap, m.GetGrades(), "grades", &v),
		Courses:                 fromProtoEnums(courseMap, m.GetCourses(), "courses", &v),
		Class:                   fromProtoEnums(classMap, m.GetClasses(), "classes", &v),
		Classification:          fromProtoEnums(classificationMap, m.GetClassifications(), "classifications", &v),
		Semester:                fromProtoEnums(semesterMap, m.GetSemesters(), "semesters", &v),
		RequirementType:         fromProtoEnums(requirementTypeMap, m.GetRequirementTypes(), "requirement_types", &v),
		CulturalSubjectCategory: fromProtoEnums(culturalCategoryMap, m.GetCulturalSubjectCategories(), "cultural_subject_categories", &v),
	}
	if err := v.err(); err != nil {
		return nil, err
	}
	subjects, err := h.service.GetSubjects(query)
	if err != nil {
		return nil, toConnectError(ctx, "ListSubjects", err)
	}
	page, next, err := paginate(m, subjects)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListSubjectsResponse{
		Subjects:      mapSlice(page, toProtoSubject),
		NextPageToken: next,
	}), nil
}

// GetSubject は科目の詳細を返す。
func (h *SubjectHandler) GetSubject(ctx context.Context, req *connect.Request[appv1.GetSubjectRequest]) (*connect.Response[appv1.GetSubjectResponse], error) {
	if req.Msg.GetId() == "" {
		return nil, invalidArgument("INVALID_ARGUMENT", "id", "is required")
	}
	subject, err := h.service.GetSubject(req.Msg.GetId())
	if err != nil {
		return nil, toConnectError(ctx, "GetSubject", err)
	}
	return connect.NewResponse(&appv1.GetSubjectResponse{Subject: toProtoSubjectDetail(*subject)}), nil
}
