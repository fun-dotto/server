package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/repository"
)

func (h *Handler) FacultyRoomsV1Create(ctx context.Context, request api.FacultyRoomsV1CreateRequestObject) (api.FacultyRoomsV1CreateResponseObject, error) {
	domainFacultyRoom := toDomainFacultyRoomFromRequest(*request.Body)
	created, err := h.facultyRoomSvc.Create(ctx, domainFacultyRoom)
	if err != nil {
		if errors.Is(err, repository.ErrFacultyRoomAlreadyExists) {
			return api.FacultyRoomsV1Create409Response{}, nil
		}
		return nil, err
	}
	return api.FacultyRoomsV1Create201JSONResponse{FacultyRoom: facultyRoomToAPI(created)}, nil
}
