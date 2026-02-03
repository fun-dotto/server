package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/announcement-api/generated"
)

func (h *Handler) AnnouncementsV1List(ctx context.Context, request api.AnnouncementsV1ListRequestObject) (api.AnnouncementsV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
