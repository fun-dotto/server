package external

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// ToDomainReservation は外部APIのReservationをDomainに変換する
func ToDomainReservation(m academic_api.Reservation) domain.Reservation {
	return domain.Reservation{
		ID: m.Id,
		Room: domain.Room{
			ID:    m.Room.Id,
			Name:  m.Room.Name,
			Floor: domain.Floor(m.Room.Floor),
		},
		StartAt: m.StartAt,
		EndAt:   m.EndAt,
		Title:   m.Title,
	}
}

// ToExternalReservationQuery はDomainのReservationQueryを外部APIのパラメータに変換する
func ToExternalReservationQuery(q domain.ReservationQuery) *academic_api.ReservationsV1ListParams {
	params := &academic_api.ReservationsV1ListParams{
		From:  q.From,
		Until: q.Until,
	}
	if len(q.RoomIDs) > 0 {
		ids := q.RoomIDs
		params.RoomIds = &ids
	}
	return params
}
