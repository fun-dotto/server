package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) FacultyRoomsV1List(_ context.Context, _ api.FacultyRoomsV1ListRequestObject) (api.FacultyRoomsV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
