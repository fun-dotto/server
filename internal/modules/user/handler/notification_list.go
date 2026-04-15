package handler

import (
	"context"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1List(_ context.Context, request api.NotificationV1ListRequestObject) (api.NotificationV1ListResponseObject, error) {
	filter := toDomainNotificationListFilter(request.Params)

	notifications, err := h.notificationService.ListNotifications(filter)
	if err != nil {
		return nil, err
	}

	return api.NotificationV1List200JSONResponse{
		Notifications: toAPINotifications(notifications),
	}, nil
}
