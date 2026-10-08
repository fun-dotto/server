package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
)

type CalendarDateHandler struct {
	c crud[model.CalendarDate]
}

func NewCalendarDateHandler(repo Repository[model.CalendarDate]) *CalendarDateHandler {
	return &CalendarDateHandler{c: crud[model.CalendarDate]{
		name: "calendar_date",
		repo: repo,
		keys: func(m model.CalendarDate) []keyField {
			return []keyField{{"service_id", m.ServiceID}, {"date", m.Date.Format(dateLayout)}}
		},
	}}
}

// calendarDateExceptionTypeEnum は GTFS の exception_type（1 / 2）との対応。
var calendarDateExceptionTypeEnum = map[adminv1.CalendarDateExceptionType]int{
	adminv1.CalendarDateExceptionType_CALENDAR_DATE_EXCEPTION_TYPE_ADDED:   1,
	adminv1.CalendarDateExceptionType_CALENDAR_DATE_EXCEPTION_TYPE_REMOVED: 2,
}

type calendarDateInput interface {
	GetServiceId() string
	GetDate() *date.Date
	GetExceptionType() adminv1.CalendarDateExceptionType
}

func calendarDateFromProto(in calendarDateInput, s fieldSet) (model.CalendarDate, error) {
	var v violations
	var m model.CalendarDate
	if s.has("service_id") {
		m.ServiceID = requireString(&v, "service_id", in.GetServiceId())
	}
	if s.has("date") {
		m.Date = requireDate(&v, "date", in.GetDate())
	}
	if s.has("exception_type") {
		t, ok := calendarDateExceptionTypeEnum[in.GetExceptionType()]
		if !ok {
			v.add("exception_type", "must be a valid enum value other than UNSPECIFIED")
		}
		m.ExceptionType = t
	}
	return m, v.err()
}

func calendarDateToProto(m model.CalendarDate) *adminv1.CalendarDate {
	d := &adminv1.CalendarDate{ServiceId: m.ServiceID, Date: toDate(m.Date)}
	for e, t := range calendarDateExceptionTypeEnum {
		if t == m.ExceptionType {
			d.ExceptionType = e
		}
	}
	return d
}

func (h *CalendarDateHandler) ListCalendarDates(ctx context.Context, req *connect.Request[adminv1.ListCalendarDatesRequest]) (*connect.Response[adminv1.ListCalendarDatesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListCalendarDatesResponse{NextPageToken: next}
	for _, m := range items {
		res.CalendarDates = append(res.CalendarDates, calendarDateToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *CalendarDateHandler) GetCalendarDate(ctx context.Context, req *connect.Request[adminv1.GetCalendarDateRequest]) (*connect.Response[adminv1.GetCalendarDateResponse], error) {
	serviceID, err := requireKeyString("service_id", req.Msg.GetServiceId())
	if err != nil {
		return nil, err
	}
	date, err := dateKey("date", req.Msg.GetDate())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"service_id", serviceID}, {"date", date}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetCalendarDateResponse{CalendarDate: calendarDateToProto(m)}), nil
}

func (h *CalendarDateHandler) CreateCalendarDate(ctx context.Context, req *connect.Request[adminv1.CreateCalendarDateRequest]) (*connect.Response[adminv1.CreateCalendarDateResponse], error) {
	m, err := calendarDateFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateCalendarDateResponse{CalendarDate: calendarDateToProto(created)}), nil
}

func (h *CalendarDateHandler) UpdateCalendarDate(ctx context.Context, req *connect.Request[adminv1.UpdateCalendarDateRequest]) (*connect.Response[adminv1.UpdateCalendarDateResponse], error) {
	serviceID, err := requireKeyString("service_id", req.Msg.GetServiceId())
	if err != nil {
		return nil, err
	}
	date, err := dateKey("date", req.Msg.GetDate())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "service_id", "date")
	if err != nil {
		return nil, err
	}
	m, err := calendarDateFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"service_id", serviceID}, {"date", date}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateCalendarDateResponse{CalendarDate: calendarDateToProto(updated)}), nil
}

func (h *CalendarDateHandler) DeleteCalendarDate(ctx context.Context, req *connect.Request[adminv1.DeleteCalendarDateRequest]) (*connect.Response[adminv1.DeleteCalendarDateResponse], error) {
	serviceID, err := requireKeyString("service_id", req.Msg.GetServiceId())
	if err != nil {
		return nil, err
	}
	date, err := dateKey("date", req.Msg.GetDate())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"service_id", serviceID}, {"date", date}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteCalendarDateResponse{}), nil
}
