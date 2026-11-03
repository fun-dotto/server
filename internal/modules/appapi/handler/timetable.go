package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
)

// requiredSemesters は必須の semesters を domain の値に変換する。
func requiredSemesters(values []appv1.SubjectSemester, v *violations) []domain.CourseSemester {
	if len(values) == 0 {
		v.add("semesters", "is required")
		return nil
	}
	return fromProtoEnums(semesterMap, values, "semesters", v)
}

func optionalYear(y *int32) *int {
	if y == nil {
		return nil
	}
	year := int(*y)
	return &year
}

// TimetableItemHandler は TimetableItemService を実装する。
type TimetableItemHandler struct {
	appv1connect.UnimplementedTimetableItemServiceHandler
	service AcademicService
}

// NewTimetableItemHandler は TimetableItemHandler を返す。
func NewTimetableItemHandler(s AcademicService) *TimetableItemHandler {
	return &TimetableItemHandler{service: s}
}

// ListTimetableItems は時間割の一覧を返す。
func (h *TimetableItemHandler) ListTimetableItems(ctx context.Context, req *connect.Request[appv1.ListTimetableItemsRequest]) (*connect.Response[appv1.ListTimetableItemsResponse], error) {
	var v violations
	semesters := requiredSemesters(req.Msg.GetSemesters(), &v)
	if err := v.err(); err != nil {
		return nil, err
	}
	items, err := h.service.GetTimetableItems(domain.TimetableItemQuery{Semesters: semesters, Year: optionalYear(req.Msg.Year)})
	if err != nil {
		return nil, toConnectError(ctx, "ListTimetableItems", err)
	}
	page, next, err := paginate(req.Msg, items)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListTimetableItemsResponse{
		TimetableItems: mapSlice(page, toProtoTimetableItem),
		NextPageToken:  next,
	}), nil
}

func toProtoTimetableItem(item domain.TimetableItem) *appv1.TimetableItem {
	p := &appv1.TimetableItem{
		Id:      item.ID,
		Subject: toProtoSubject(item.Subject),
		Rooms:   toProtoRooms(item.Rooms),
	}
	if item.Slot != nil {
		p.Slot = &appv1.TimetableItemSlot{
			DayOfWeek: dayOfWeekMap.toProto[item.Slot.DayOfWeek],
			Period:    periodMap.toProto[item.Slot.Period],
		}
	}
	return p
}

// CourseRegistrationHandler は CourseRegistrationService を実装する。
type CourseRegistrationHandler struct {
	appv1connect.UnimplementedCourseRegistrationServiceHandler
	service AcademicService
}

// NewCourseRegistrationHandler は CourseRegistrationHandler を返す。
func NewCourseRegistrationHandler(s AcademicService) *CourseRegistrationHandler {
	return &CourseRegistrationHandler{service: s}
}

// ListCourseRegistrations はログイン中のユーザーの履修情報の一覧を返す。
func (h *CourseRegistrationHandler) ListCourseRegistrations(ctx context.Context, req *connect.Request[appv1.ListCourseRegistrationsRequest]) (*connect.Response[appv1.ListCourseRegistrationsResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	var v violations
	semesters := requiredSemesters(req.Msg.GetSemesters(), &v)
	if err := v.err(); err != nil {
		return nil, err
	}
	registrations, err := h.service.GetCourseRegistrations(user.ID, semesters, optionalYear(req.Msg.Year))
	if err != nil {
		return nil, toConnectError(ctx, "ListCourseRegistrations", err)
	}
	page, next, err := paginate(req.Msg, registrations)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListCourseRegistrationsResponse{
		CourseRegistrations: mapSlice(page, toProtoCourseRegistration),
		NextPageToken:       next,
	}), nil
}

// CreateCourseRegistration はログイン中のユーザーの履修情報を作成する。
func (h *CourseRegistrationHandler) CreateCourseRegistration(ctx context.Context, req *connect.Request[appv1.CreateCourseRegistrationRequest]) (*connect.Response[appv1.CreateCourseRegistrationResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	if req.Msg.GetSubjectId() == "" {
		return nil, invalidArgument("INVALID_ARGUMENT", "subject_id", "is required")
	}
	registration, err := h.service.CreateCourseRegistration(user.ID, req.Msg.GetSubjectId())
	if err != nil {
		return nil, toConnectError(ctx, "CreateCourseRegistration", err)
	}
	return connect.NewResponse(&appv1.CreateCourseRegistrationResponse{
		CourseRegistration: toProtoCourseRegistration(*registration),
	}), nil
}

// DeleteCourseRegistration は履修情報を削除する。
func (h *CourseRegistrationHandler) DeleteCourseRegistration(ctx context.Context, req *connect.Request[appv1.DeleteCourseRegistrationRequest]) (*connect.Response[appv1.DeleteCourseRegistrationResponse], error) {
	if _, ok := auth.UserFromContext(ctx); !ok {
		return nil, unauthenticated()
	}
	if req.Msg.GetId() == "" {
		return nil, invalidArgument("INVALID_ARGUMENT", "id", "is required")
	}
	if err := h.service.DeleteCourseRegistration(req.Msg.GetId()); err != nil {
		return nil, toConnectError(ctx, "DeleteCourseRegistration", err)
	}
	return connect.NewResponse(&appv1.DeleteCourseRegistrationResponse{}), nil
}

func toProtoCourseRegistration(r domain.CourseRegistration) *appv1.CourseRegistration {
	return &appv1.CourseRegistration{Id: r.ID, Subject: toProtoSubject(r.Subject)}
}

// PersonalCalendarItemHandler は PersonalCalendarItemService を実装する。
type PersonalCalendarItemHandler struct {
	appv1connect.UnimplementedPersonalCalendarItemServiceHandler
	service AcademicService
}

// NewPersonalCalendarItemHandler は PersonalCalendarItemHandler を返す。
func NewPersonalCalendarItemHandler(s AcademicService) *PersonalCalendarItemHandler {
	return &PersonalCalendarItemHandler{service: s}
}

// ListPersonalCalendarItems はログイン中のユーザーの個人カレンダーアイテムの一覧を返す。
func (h *PersonalCalendarItemHandler) ListPersonalCalendarItems(ctx context.Context, req *connect.Request[appv1.ListPersonalCalendarItemsRequest]) (*connect.Response[appv1.ListPersonalCalendarItemsResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	var v violations
	if len(req.Msg.GetDates()) == 0 {
		v.add("dates", "is required")
	}
	dates := make([]time.Time, 0, len(req.Msg.GetDates()))
	for _, d := range req.Msg.GetDates() {
		t, ok := fromProtoDate(d)
		if !ok {
			v.add("dates", "must contain only valid dates")
			continue
		}
		dates = append(dates, t)
	}
	if err := v.err(); err != nil {
		return nil, err
	}
	items, err := h.service.GetPersonalCalendarItems(user.ID, dates)
	if err != nil {
		return nil, toConnectError(ctx, "ListPersonalCalendarItems", err)
	}
	page, next, err := paginate(req.Msg, items)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListPersonalCalendarItemsResponse{
		PersonalCalendarItems: mapSlice(page, func(item domain.PersonalCalendarItem) *appv1.PersonalCalendarItem {
			return &appv1.PersonalCalendarItem{
				Date:    toProtoDate(item.Date),
				Period:  periodMap.toProto[item.Period],
				Subject: toProtoSubject(item.Subject),
				Rooms:   toProtoRooms(item.Rooms),
				State:   calendarStateMap.toProto[item.Status],
			}
		}),
		NextPageToken: next,
	}), nil
}
