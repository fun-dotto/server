package handler

import (
	"context"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) NotificationV1Dispatch(ctx context.Context, request api.NotificationV1DispatchRequestObject) (api.NotificationV1DispatchResponseObject, error) {
	if request.Body == nil || len(request.Body.NotificationIds) == 0 {
		return api.NotificationV1Dispatch400Response{}, nil
	}

	notifications, err := h.notificationService.DispatchNotifications(ctx, request.Body.NotificationIds)
	if err != nil {
		return nil, err
	}

	return api.NotificationV1Dispatch200JSONResponse{
		Notifications: toAPINotifications(notifications),
	}, nil
}
