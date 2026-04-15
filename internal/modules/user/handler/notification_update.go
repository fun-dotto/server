package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/user-api/generated"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (h *Handler) NotificationV1Update(_ context.Context, request api.NotificationV1UpdateRequestObject) (api.NotificationV1UpdateResponseObject, error) {
	notification := toDomainNotification(request.Id, *request.Body)

	updated, err := h.notificationService.UpdateNotification(notification)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.NotificationV1Update404Response{}, nil
		}
		return nil, err
	}

	return api.NotificationV1Update200JSONResponse{
		Notification: toAPINotification(updated),
	}, nil
}
