package handler

import (
	"context"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) FCMTokenV1List(_ context.Context, request api.FCMTokenV1ListRequestObject) (api.FCMTokenV1ListResponseObject, error) {
	filter := toDomainFCMTokenListFilter(request.Params)

	tokens, err := h.fcmTokenService.ListFCMTokens(filter)
	if err != nil {
		return nil, err
	}

	return api.FCMTokenV1List200JSONResponse{
		FcmTokens: toAPIFCMTokens(tokens),
	}, nil
}
