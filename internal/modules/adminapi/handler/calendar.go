package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
)

type CalendarHandler struct {
	c crud[model.Calendar]
}

func NewCalendarHandler(repo Repository[model.Calendar]) *CalendarHandler {
	return &CalendarHandler{c: crud[model.Calendar]{
		name: "calendar",
		repo: repo,
		keys: func(m model.Calendar) []keyField { return []keyField{{"service_id", m.ServiceID}} },
	}}
}

type calendarInput interface {
	GetServiceId() string
	GetMonday() bool
	GetTuesday() bool
	GetWednesday() bool
	GetThursday() bool
	GetFriday() bool
	GetSaturday() bool
	GetSunday() bool
	GetStartDate() *date.Date
	GetEndDate() *date.Date
}

func calendarFromProto(in calendarInput, s fieldSet) (model.Calendar, error) {
	var v violations
	m := model.Calendar{
		Monday:    in.GetMonday(),
		Tuesday:   in.GetTuesday(),
		Wednesday: in.GetWednesday(),
		Thursday:  in.GetThursday(),
		Friday:    in.GetFriday(),
		Saturday:  in.GetSaturday(),
		Sunday:    in.GetSunday(),
	}
	if s.has("service_id") {
		m.ServiceID = requireString(&v, "service_id", in.GetServiceId())
	}
	if s.has("start_date") {
		m.StartDate = requireDate(&v, "start_date", in.GetStartDate())
	}
	if s.has("end_date") {
		m.EndDate = requireDate(&v, "end_date", in.GetEndDate())
	}
	if s.has("start_date") && s.has("end_date") && len(v) == 0 && m.EndDate.Before(m.StartDate) {
		v.add("end_date", "must not be before start_date")
	}
	return m, v.err()
}

func calendarToProto(m model.Calendar) *adminv1.Calendar {
	return &adminv1.Calendar{
		ServiceId: m.ServiceID,
		Monday:    m.Monday,
		Tuesday:   m.Tuesday,
		Wednesday: m.Wednesday,
		Thursday:  m.Thursday,
		Friday:    m.Friday,
		Saturday:  m.Saturday,
		Sunday:    m.Sunday,
		StartDate: toDate(m.StartDate),
		EndDate:   toDate(m.EndDate),
	}
}

func (h *CalendarHandler) ListCalendars(ctx context.Context, req *connect.Request[adminv1.ListCalendarsRequest]) (*connect.Response[adminv1.ListCalendarsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListCalendarsResponse{NextPageToken: next}
	for _, m := range items {
		res.Calendars = append(res.Calendars, calendarToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *CalendarHandler) GetCalendar(ctx context.Context, req *connect.Request[adminv1.GetCalendarRequest]) (*connect.Response[adminv1.GetCalendarResponse], error) {
	serviceID, err := requireKeyString("service_id", req.Msg.GetServiceId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"service_id", serviceID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetCalendarResponse{Calendar: calendarToProto(m)}), nil
}

func (h *CalendarHandler) CreateCalendar(ctx context.Context, req *connect.Request[adminv1.CreateCalendarRequest]) (*connect.Response[adminv1.CreateCalendarResponse], error) {
	m, err := calendarFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateCalendarResponse{Calendar: calendarToProto(created)}), nil
}

func (h *CalendarHandler) UpdateCalendar(ctx context.Context, req *connect.Request[adminv1.UpdateCalendarRequest]) (*connect.Response[adminv1.UpdateCalendarResponse], error) {
	serviceID, err := requireKeyString("service_id", req.Msg.GetServiceId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "service_id")
	if err != nil {
		return nil, err
	}
	m, err := calendarFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"service_id", serviceID}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateCalendarResponse{Calendar: calendarToProto(updated)}), nil
}

func (h *CalendarHandler) DeleteCalendar(ctx context.Context, req *connect.Request[adminv1.DeleteCalendarRequest]) (*connect.Response[adminv1.DeleteCalendarResponse], error) {
	serviceID, err := requireKeyString("service_id", req.Msg.GetServiceId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"service_id", serviceID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteCalendarResponse{}), nil
}
