package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultyRoomsV1Delete(_ context.Context, _ api.FacultyRoomsV1DeleteRequestObject) (api.FacultyRoomsV1DeleteResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
