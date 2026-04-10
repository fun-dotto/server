package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1Delete(_ context.Context, _ api.NotificationV1DeleteRequestObject) (api.NotificationV1DeleteResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
