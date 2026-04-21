package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1Dispatch(_ context.Context, _ api.NotificationV1DispatchRequestObject) (api.NotificationV1DispatchResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
