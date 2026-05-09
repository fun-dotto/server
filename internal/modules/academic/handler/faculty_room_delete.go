package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/server/gen/academic"
	"gorm.io/gorm"
)

func (h *Handler) FacultyRoomsV1Delete(ctx context.Context, request api.FacultyRoomsV1DeleteRequestObject) (api.FacultyRoomsV1DeleteResponseObject, error) {
	if err := h.facultyRoomSvc.Delete(ctx, request.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.FacultyRoomsV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.FacultyRoomsV1Delete204Response{}, nil
}
