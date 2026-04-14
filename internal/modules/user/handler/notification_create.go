package handler

import (
	"context"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1Create(_ context.Context, request api.NotificationV1CreateRequestObject) (api.NotificationV1CreateResponseObject, error) {
	notification := toDomainNotification("", *request.Body)

	created, err := h.notificationService.CreateNotification(notification)
	if err != nil {
		return nil, err
	}

	return api.NotificationV1Create201JSONResponse{
		Notification: toAPINotification(created),
	}, nil
}
