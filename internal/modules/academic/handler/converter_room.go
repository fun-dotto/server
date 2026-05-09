package handler

import (
	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func buildRoomListFilter(params api.RoomsV1ListParams) domain.RoomListFilter {
	filter := domain.RoomListFilter{}
	if params.Floors != nil {
		floors := make([]domain.Floor, len(*params.Floors))
		for i, f := range *params.Floors {
			floors[i] = domain.Floor(f)
		}
		filter.Floors = floors
	}
	return filter
}

func roomToAPI(room domain.Room) api.Room {
	return api.Room{
		Id:    room.ID,
		Name:  room.Name,
		Floor: api.DottoFoundationV1Floor(room.Floor),
	}
}

func roomsToAPI(rooms []domain.Room) []api.Room {
	result := make([]api.Room, len(rooms))
	for i, room := range rooms {
		result[i] = roomToAPI(room)
	}
	return result
}

func toDomainRoomFromRequest(id string, req api.RoomRequest) domain.Room {
	return domain.Room{
		ID:    id,
		Name:  req.Name,
		Floor: domain.Floor(req.Floor),
	}
}
