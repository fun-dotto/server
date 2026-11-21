package handler

import (
	"context"

	"connectrpc.com/connect"

	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/fun-dotto/server/internal/shared/rpcerror"
)

// ListReservations は教室の予約の一覧を返す。検索対象期間に一部でも重複する予約を返す。
func (h *Handler) ListReservations(ctx context.Context, req *connect.Request[appv1.ListReservationsRequest]) (*connect.Response[appv1.ListReservationsResponse], error) {
	m := req.Msg
	var v rpcerror.Violations
	roomIDs := parseUUIDs(m.GetRoomIds(), "room_ids", &v)
	query := h.db.WithContext(ctx).Preload("Room")
	if len(roomIDs) > 0 {
		query = query.Where("room_id IN ?", roomIDs)
	}
	if m.GetFromTime() != nil {
		if err := m.GetFromTime().CheckValid(); err != nil {
			v.Add("from_time", "must be a valid timestamp")
		}
		query = query.Where("end_time > ?", m.GetFromTime().AsTime())
	}
	if m.GetUntilTime() != nil {
		if err := m.GetUntilTime().CheckValid(); err != nil {
			v.Add("until_time", "must be a valid timestamp")
		}
		query = query.Where("start_time < ?", m.GetUntilTime().AsTime())
	}
	if err := v.Err(errorDomain); err != nil {
		return nil, err
	}
	var records []model.RoomReservation
	if err := query.Order("start_time ASC, id ASC").Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListReservations", err)
	}
	reservations := mapSlice(records, func(r model.RoomReservation) *appv1.Reservation {
		return &appv1.Reservation{
			Id:      r.ID.String(),
			Room:    toProtoRoom(r.Room),
			StartAt: toTimestamp(r.StartTime),
			EndAt:   toTimestamp(r.EndTime),
			Title:   r.Title,
		}
	})
	rooms := mapSlice(reservations, (*appv1.Reservation).GetRoom)
	if err := h.attachRoomFaculties(ctx, rooms...); err != nil {
		return nil, internalError(ctx, "ListReservations", err)
	}
	return connect.NewResponse(&appv1.ListReservationsResponse{Reservations: reservations}), nil
}
