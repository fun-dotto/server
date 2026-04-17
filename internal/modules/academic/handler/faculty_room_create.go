package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultyRoomsV1Create(_ context.Context, _ api.FacultyRoomsV1CreateRequestObject) (api.FacultyRoomsV1CreateResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
