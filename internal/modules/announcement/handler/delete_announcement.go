package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/announcement-api/generated"
)

func (h *Handler) AnnouncementsV1Delete(ctx context.Context, request api.AnnouncementsV1DeleteRequestObject) (api.AnnouncementsV1DeleteResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
