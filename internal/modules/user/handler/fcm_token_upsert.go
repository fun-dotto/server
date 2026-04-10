package handler

import (
	"context"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) FCMTokenV1Upsert(_ context.Context, request api.FCMTokenV1UpsertRequestObject) (api.FCMTokenV1UpsertResponseObject, error) {
	fcmToken, err := h.fcmTokenService.UpsertFCMToken(toDomainFCMToken(*request.Body))
	if err != nil {
		return nil, err
	}

	return api.FCMTokenV1Upsert200JSONResponse{
		FcmToken: toAPIFCMToken(fcmToken),
	}, nil
}
