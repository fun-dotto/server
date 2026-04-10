package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1Update(_ context.Context, _ api.NotificationV1UpdateRequestObject) (api.NotificationV1UpdateResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
