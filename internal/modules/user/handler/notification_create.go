package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1Create(_ context.Context, _ api.NotificationV1CreateRequestObject) (api.NotificationV1CreateResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
