package handler

import (
	"context"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) FCMTokenV1List(ctx context.Context, request api.FCMTokenV1ListRequestObject) (api.FCMTokenV1ListResponseObject, error) {
	filter := toDomainFCMTokenListFilter(request.Params)

	tokens, err := h.fcmTokenService.ListFCMTokens(ctx, filter)
	if err != nil {
		return nil, err
	}

	return api.FCMTokenV1List200JSONResponse{
		FcmTokens: toAPIFCMTokens(tokens),
	}, nil
}
