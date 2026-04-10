package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1List(_ context.Context, _ api.NotificationV1ListRequestObject) (api.NotificationV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
