package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// ReservationHandler は ReservationService を実装する。
type ReservationHandler struct {
	appv1connect.UnimplementedReservationServiceHandler
	service AcademicService
}

// NewReservationHandler は ReservationHandler を返す。
func NewReservationHandler(s AcademicService) *ReservationHandler {
	return &ReservationHandler{service: s}
}

// ListReservations は教室の予約の一覧を返す。
func (h *ReservationHandler) ListReservations(ctx context.Context, req *connect.Request[appv1.ListReservationsRequest]) (*connect.Response[appv1.ListReservationsResponse], error) {
	m := req.Msg
	var v violations
	query := domain.ReservationQuery{RoomIDs: m.GetRoomIds()}
	if m.FromTime != nil {
		if err := m.FromTime.CheckValid(); err != nil {
			v.add("from_time", "must be a valid timestamp")
		}
		t := m.FromTime.AsTime()
		query.From = &t
	}
	if m.UntilTime != nil {
		if err := m.UntilTime.CheckValid(); err != nil {
			v.add("until_time", "must be a valid timestamp")
		}
		t := m.UntilTime.AsTime()
		query.Until = &t
	}
	if err := v.err(); err != nil {
		return nil, err
	}
	reservations, err := h.service.GetReservations(query)
	if err != nil {
		return nil, toConnectError(ctx, "ListReservations", err)
	}
	page, next, err := paginate(m, reservations)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListReservationsResponse{
		Reservations: mapSlice(page, func(r domain.Reservation) *appv1.Reservation {
			return &appv1.Reservation{
				Id:      r.ID,
				Room:    toProtoRoom(r.Room),
				StartAt: toTimestamp(r.StartAt),
				EndAt:   toTimestamp(r.EndAt),
				Title:   r.Title,
			}
		}),
		NextPageToken: next,
	}), nil
}
