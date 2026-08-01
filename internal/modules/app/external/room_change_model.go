package external

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDomainRoomChange は外部APIのRoomChangeをDomainに変換する
func ToDomainRoomChange(m academic_api.RoomChange) domain.RoomChange {
	return domain.RoomChange{
		ID:      m.Id,
		Date:    m.Date.Time,
		Period:  domain.Period(m.Period),
		Subject: ToDomainSubjectSummary(m.Subject),
		OriginalRoom: domain.Room{
			ID:    m.OriginalRoom.Id,
			Name:  m.OriginalRoom.Name,
			Floor: domain.Floor(m.OriginalRoom.Floor),
		},
		NewRoom: domain.Room{
			ID:    m.NewRoom.Id,
			Name:  m.NewRoom.Name,
			Floor: domain.Floor(m.NewRoom.Floor),
		},
	}
}

// ToExternalRoomChangeQuery はDomainのRoomChangeQueryを外部APIのパラメータに変換する
func ToExternalRoomChangeQuery(q domain.RoomChangeQuery) *academic_api.RoomChangesV1ListParams {
	params := &academic_api.RoomChangesV1ListParams{
		SubjectIds: q.SubjectIDs,
	}
	if q.From != nil {
		d := openapi_types.Date{Time: *q.From}
		params.From = &d
	}
	if q.Until != nil {
		d := openapi_types.Date{Time: *q.Until}
		params.Until = &d
	}
	return params
}
