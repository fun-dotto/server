package handler

import (
	"time"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

var nowFunc = time.Now

func buildFacultyRoomListFilter(params api.FacultyRoomsV1ListParams) domain.FacultyRoomListFilter {
	filter := domain.FacultyRoomListFilter{}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		currentYear := domain.CurrentAcademicYear(nowFunc())
		filter.Year = &currentYear
	}
	return filter
}

func facultyRoomToAPI(fr domain.FacultyRoom) api.FacultyRoom {
	return api.FacultyRoom{
		Id:      fr.ID,
		Faculty: facultyToAPI(fr.Faculty),
		Room:    roomToAPI(fr.Room),
		Year:    fr.Year,
	}
}

func facultyRoomsToAPI(facultyRooms []domain.FacultyRoom) []api.FacultyRoom {
	result := make([]api.FacultyRoom, len(facultyRooms))
	for i, fr := range facultyRooms {
		result[i] = facultyRoomToAPI(fr)
	}
	return result
}

func toDomainFacultyRoomFromRequest(req api.FacultyRoomRequest) domain.FacultyRoom {
	return domain.FacultyRoom{
		Faculty: domain.Faculty{ID: req.FacultyId},
		Room:    domain.Room{ID: req.RoomId},
		Year:    req.Year,
	}
}
