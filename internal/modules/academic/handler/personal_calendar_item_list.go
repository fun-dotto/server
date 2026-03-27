package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) PersonalCalendarItemsV1List(
	ctx context.Context,
	request api.PersonalCalendarItemsV1ListRequestObject) (api.PersonalCalendarItemsV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
