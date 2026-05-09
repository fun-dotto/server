package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultyRoomsV1List(ctx context.Context, request api.FacultyRoomsV1ListRequestObject) (api.FacultyRoomsV1ListResponseObject, error) {
	filter := buildFacultyRoomListFilter(request.Params)

	facultyRooms, err := h.facultyRoomSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return api.FacultyRoomsV1List200JSONResponse{FacultyRooms: facultyRoomsToAPI(facultyRooms)}, nil
}
