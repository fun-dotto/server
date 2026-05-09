package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) CancelledClassesV1Create(ctx context.Context, request api.CancelledClassesV1CreateRequestObject) (api.CancelledClassesV1CreateResponseObject, error) {
	domainCC := toDomainCancelledClassFromRequest(*request.Body)
	created, err := h.cancelledClassSvc.Create(ctx, domainCC)
	if err != nil {
		return nil, err
	}
	res, err := cancelledClassToAPI(created)
	if err != nil {
		return nil, err
	}
	return api.CancelledClassesV1Create201JSONResponse{CancelledClass: res}, nil
}
