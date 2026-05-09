package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) CancelledClassesV1List(ctx context.Context, request api.CancelledClassesV1ListRequestObject) (api.CancelledClassesV1ListResponseObject, error) {
	filter := buildCancelledClassListFilter(request.Params)

	cancelledClasses, err := h.cancelledClassSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	res, err := cancelledClassesToAPI(cancelledClasses)
	if err != nil {
		return nil, err
	}
	return api.CancelledClassesV1List200JSONResponse{CancelledClasses: res}, nil
}
