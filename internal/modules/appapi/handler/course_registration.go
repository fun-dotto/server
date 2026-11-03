package handler

import (
	"context"
	"strings"
	"uuid"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

// courseRegistrationIDSep は履修情報ID（ユーザーIDと科目IDの複合キー）の区切り文字。
const courseRegistrationIDSep = "|"

func encodeCourseRegistrationID(userID string, subjectID uuid.UUID) string {
	return userID + courseRegistrationIDSep + subjectID.String()
}

func decodeCourseRegistrationID(id string) (userID string, subjectID uuid.UUID, ok bool) {
	parts := strings.SplitN(id, courseRegistrationIDSep, 2)
	if len(parts) != 2 || parts[0] == "" {
		return "", uuid.Nil(), false
	}
	sid, err := uuid.Parse(parts[1])
	if err != nil {
		return "", uuid.Nil(), false
	}
	return parts[0], sid, true
}

func toProtoCourseRegistration(r model.CourseRegistration) *appv1.CourseRegistration {
	return &appv1.CourseRegistration{
		Id:      encodeCourseRegistrationID(r.UserID, r.SubjectID),
		Subject: toProtoSubject(r.Subject),
	}
}

// ListCourseRegistrations はログイン中のユーザーの履修情報の一覧を返す。
func (h *Handler) ListCourseRegistrations(ctx context.Context, req *connect.Request[appv1.ListCourseRegistrationsRequest]) (*connect.Response[appv1.ListCourseRegistrationsResponse], error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var v violations
	semesters := requiredSemesters(req.Msg.GetSemesters(), &v)
	if err := v.err(); err != nil {
		return nil, err
	}
	var records []model.CourseRegistration
	if err := preload(h.db.WithContext(ctx), subjectPreloads("Subject.")...).
		Joins("JOIN subjects ON subjects.id = course_registrations.subject_id").
		Where("course_registrations.user_id = ?", userID).
		Where("subjects.year = ?", h.defaultYear(req.Msg.Year)).
		Where("subjects.semester IN ?", semesters).
		Order("subjects.syllabus_id ASC").
		Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListCourseRegistrations", err)
	}
	return connect.NewResponse(&appv1.ListCourseRegistrationsResponse{
		CourseRegistrations: mapSlice(records, toProtoCourseRegistration),
	}), nil
}

// CreateCourseRegistration はログイン中のユーザーの履修情報を作成する。
func (h *Handler) CreateCourseRegistration(ctx context.Context, req *connect.Request[appv1.CreateCourseRegistrationRequest]) (*connect.Response[appv1.CreateCourseRegistrationResponse], error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	subjectID, err := uuid.Parse(req.Msg.GetSubjectId())
	if err != nil {
		return nil, invalidArgument("subject_id", "must be a valid subject ID")
	}
	metadata := map[string]string{"subjectId": subjectID.String()}

	db := h.db.WithContext(ctx)
	var count int64
	if err := db.Model(&model.Subject{}).Where("id = ?", subjectID.String()).Count(&count).Error; err != nil {
		return nil, internalError(ctx, "CreateCourseRegistration", err)
	}
	if count == 0 {
		return nil, notFound("subject", metadata)
	}

	record := model.CourseRegistration{UserID: userID, SubjectID: subjectID}
	if err := db.Omit("User", "Subject").Create(&record).Error; err != nil {
		switch postgresErrorCode(err) {
		case pgUniqueViolation:
			return nil, newError(connect.CodeAlreadyExists, "COURSE_REGISTRATION_ALREADY_EXISTS",
				"course registration already exists for the subject", metadata)
		case pgForeignKeyViolation:
			return nil, userNotFound()
		}
		return nil, internalError(ctx, "CreateCourseRegistration", err)
	}

	var created model.CourseRegistration
	if err := preload(db, subjectPreloads("Subject.")...).
		First(&created, "user_id = ? AND subject_id = ?", userID, subjectID.String()).Error; err != nil {
		return nil, internalError(ctx, "CreateCourseRegistration", err)
	}
	return connect.NewResponse(&appv1.CreateCourseRegistrationResponse{
		CourseRegistration: toProtoCourseRegistration(created),
	}), nil
}

// DeleteCourseRegistration はログイン中のユーザーの履修情報を削除する。
// ログイン中のユーザーの履修情報でない場合は PERMISSION_DENIED を返す。
func (h *Handler) DeleteCourseRegistration(ctx context.Context, req *connect.Request[appv1.DeleteCourseRegistrationRequest]) (*connect.Response[appv1.DeleteCourseRegistrationResponse], error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	id := req.Msg.GetId()
	if id == "" {
		return nil, invalidArgument("id", "is required")
	}
	metadata := map[string]string{"courseRegistrationId": id}
	ownerID, subjectID, ok := decodeCourseRegistrationID(id)
	if !ok || ownerID != userID {
		return nil, permissionDenied("course_registration", metadata)
	}
	result := h.db.WithContext(ctx).
		Where("user_id = ? AND subject_id = ?", userID, subjectID.String()).
		Delete(&model.CourseRegistration{})
	if result.Error != nil {
		return nil, internalError(ctx, "DeleteCourseRegistration", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, notFound("course_registration", metadata)
	}
	return connect.NewResponse(&appv1.DeleteCourseRegistrationResponse{}), nil
}
