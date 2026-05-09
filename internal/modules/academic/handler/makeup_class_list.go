package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) MakeupClassesV1List(ctx context.Context, request api.MakeupClassesV1ListRequestObject) (api.MakeupClassesV1ListResponseObject, error) {
	filter := buildMakeupClassListFilter(request.Params)

	makeupClasses, err := h.makeupClassSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	res, err := makeupClassesToAPI(makeupClasses)
	if err != nil {
		return nil, err
	}
	return api.MakeupClassesV1List200JSONResponse{MakeupClasses: res}, nil
}
