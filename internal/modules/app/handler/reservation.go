package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

func (h *Handler) ReservationsV1List(_ context.Context, request api.ReservationsV1ListRequestObject) (api.ReservationsV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toReservationQuery(request.Params)

	reservations, err := h.academicService.GetReservations(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	apiReservations := make([]api.Reservation, len(reservations))
	for i, r := range reservations {
		apiReservations[i] = toApiReservation(r)
	}

	return api.ReservationsV1List200JSONResponse{
		Reservations: apiReservations,
	}, nil
}

func toReservationQuery(params api.ReservationsV1ListParams) domain.ReservationQuery {
	query := domain.ReservationQuery{
		From:  params.From,
		Until: params.Until,
	}
	if params.RoomIds != nil {
		query.RoomIDs = *params.RoomIds
	}
	return query
}

func toApiReservation(r domain.Reservation) api.Reservation {
	return api.Reservation{
		Id: r.ID,
		Room: api.Room{
			Id:    r.Room.ID,
			Name:  r.Room.Name,
			Floor: api.DottoFoundationV1Floor(r.Room.Floor),
		},
		StartAt: r.StartAt,
		EndAt:   r.EndAt,
		Title:   r.Title,
	}
}
