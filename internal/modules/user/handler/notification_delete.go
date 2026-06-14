package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/user-api/generated"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (h *Handler) NotificationV1Delete(ctx context.Context, request api.NotificationV1DeleteRequestObject) (api.NotificationV1DeleteResponseObject, error) {
	err := h.notificationService.DeleteNotification(ctx, request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.NotificationV1Delete404Response{}, nil
		}
		return nil, err
	}

	return api.NotificationV1Delete204Response{}, nil
}
