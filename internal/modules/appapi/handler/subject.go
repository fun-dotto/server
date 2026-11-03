package handler

import (
	"context"
	"strings"
	"uuid"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

// ListSubjects は科目を検索する。同一項目同士は OR、異なる項目同士は AND で絞り込む。
func (h *Handler) ListSubjects(ctx context.Context, req *connect.Request[appv1.ListSubjectsRequest]) (*connect.Response[appv1.ListSubjectsResponse], error) {
	m := req.Msg
	var v violations
	grades := gradeMap.fromProtoList(m.GetGrades(), "grades", &v)
	courses := courseMap.fromProtoList(m.GetCourses(), "courses", &v)
	classes := classMap.fromProtoList(m.GetClasses(), "classes", &v)
	classifications := classificationMap.fromProtoList(m.GetClassifications(), "classifications", &v)
	semesters := semesterMap.fromProtoList(m.GetSemesters(), "semesters", &v)
	requirementTypes := requirementTypeMap.fromProtoList(m.GetRequirementTypes(), "requirement_types", &v)
	culturalCategories := culturalCategoryMap.fromProtoList(m.GetCulturalSubjectCategories(), "cultural_subject_categories", &v)
	if err := v.err(); err != nil {
		return nil, err
	}

	db := h.db.WithContext(ctx)
	query := preload(db, subjectPreloads("")...)
	if q := strings.TrimSpace(m.GetQuery()); q != "" {
		query = query.Where("name ILIKE ? ESCAPE '\\'", "%"+escapeLike(q)+"%")
	}
	if len(semesters) > 0 {
		query = query.Where("semester IN ?", semesters)
	}
	if len(classifications) > 0 {
		query = query.Where("classification IN ?", classifications)
	}
	if len(culturalCategories) > 0 {
		query = query.Where("cultural_subject_category IN ?", culturalCategories)
	}
	if len(grades) > 0 || len(classes) > 0 {
		sub := db.Model(&model.SubjectEligibleAttribute{}).Select("subject_id")
		if len(grades) > 0 {
			sub = sub.Where("grade IN ?", grades)
		}
		if len(classes) > 0 {
			sub = sub.Where("class IN ?", classes)
		}
		query = query.Where("id IN (?)", sub)
	}
	if len(courses) > 0 || len(requirementTypes) > 0 {
		sub := db.Model(&model.SubjectRequirement{}).Select("subject_id")
		if len(courses) > 0 {
			sub = sub.Where("course IN ?", courses)
		}
		if len(requirementTypes) > 0 {
			sub = sub.Where("requirement_type IN ?", requirementTypes)
		}
		query = query.Where("id IN (?)", sub)
	}

	var records []model.Subject
	if err := query.Order("syllabus_id ASC").Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListSubjects", err)
	}
	return connect.NewResponse(&appv1.ListSubjectsResponse{
		Subjects: mapSlice(records, func(s model.Subject) *appv1.Subject { return toProtoSubject(&s) }),
	}), nil
}

// GetSubject は科目の詳細を返す。
func (h *Handler) GetSubject(ctx context.Context, req *connect.Request[appv1.GetSubjectRequest]) (*connect.Response[appv1.GetSubjectResponse], error) {
	id, err := uuid.Parse(req.Msg.GetId())
	if err != nil {
		return nil, invalidArgument("id", "must be a valid subject ID")
	}
	var record model.Subject
	if err := preload(h.db.WithContext(ctx), subjectDetailPreloads()...).First(&record, "id = ?", id.String()).Error; err != nil {
		if isNotFound(err) {
			return nil, notFound("subject", map[string]string{"subjectId": id.String()})
		}
		return nil, internalError(ctx, "GetSubject", err)
	}
	return connect.NewResponse(&appv1.GetSubjectResponse{Subject: toProtoSubjectDetail(&record)}), nil
}

// escapeLike は LIKE のワイルドカード文字をエスケープする。
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
