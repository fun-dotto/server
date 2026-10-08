package handler

import (
	"context"
	"uuid"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

type CourseRegistrationHandler struct {
	c crud[model.CourseRegistration]
}

func NewCourseRegistrationHandler(repo Repository[model.CourseRegistration]) *CourseRegistrationHandler {
	return &CourseRegistrationHandler{c: crud[model.CourseRegistration]{
		name: "course_registration",
		repo: repo,
		keys: func(m model.CourseRegistration) []keyField {
			return []keyField{{"user_id", m.UserID}, {"subject_id", m.SubjectID.String()}}
		},
	}}
}

type courseRegistrationInput interface {
	GetUserId() string
	GetSubjectId() string
}

func courseRegistrationFromProto(in courseRegistrationInput, s fieldSet) (model.CourseRegistration, error) {
	var v violations
	var m model.CourseRegistration
	if s.has("user_id") {
		m.UserID = requireString(&v, "user_id", in.GetUserId())
	}
	if s.has("subject_id") {
		m.SubjectID = parseUUID(&v, "subject_id", in.GetSubjectId())
	}
	return m, v.err()
}

func courseRegistrationToProto(m model.CourseRegistration) *adminv1.CourseRegistration {
	return &adminv1.CourseRegistration{
		UserId:    m.UserID,
		SubjectId: m.SubjectID.String(),
		CreatedAt: toTimestamp(m.CreatedAt),
		UpdatedAt: toTimestamp(m.UpdatedAt),
	}
}

func (h *CourseRegistrationHandler) ListCourseRegistrations(ctx context.Context, req *connect.Request[adminv1.ListCourseRegistrationsRequest]) (*connect.Response[adminv1.ListCourseRegistrationsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListCourseRegistrationsResponse{NextPageToken: next}
	for _, m := range items {
		res.CourseRegistrations = append(res.CourseRegistrations, courseRegistrationToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *CourseRegistrationHandler) GetCourseRegistration(ctx context.Context, req *connect.Request[adminv1.GetCourseRegistrationRequest]) (*connect.Response[adminv1.GetCourseRegistrationResponse], error) {
	userID, err := requireKeyString("user_id", req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	subjectID, err := parseIDField("subject_id", req.Msg.GetSubjectId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"user_id", userID}, {"subject_id", subjectID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetCourseRegistrationResponse{CourseRegistration: courseRegistrationToProto(m)}), nil
}

func (h *CourseRegistrationHandler) CreateCourseRegistration(ctx context.Context, req *connect.Request[adminv1.CreateCourseRegistrationRequest]) (*connect.Response[adminv1.CreateCourseRegistrationResponse], error) {
	m, err := courseRegistrationFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	m.Id = uuid.New()
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateCourseRegistrationResponse{CourseRegistration: courseRegistrationToProto(created)}), nil
}

func (h *CourseRegistrationHandler) DeleteCourseRegistration(ctx context.Context, req *connect.Request[adminv1.DeleteCourseRegistrationRequest]) (*connect.Response[adminv1.DeleteCourseRegistrationResponse], error) {
	userID, err := requireKeyString("user_id", req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	subjectID, err := parseIDField("subject_id", req.Msg.GetSubjectId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"user_id", userID}, {"subject_id", subjectID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteCourseRegistrationResponse{}), nil
}
