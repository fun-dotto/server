package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) MakeupClassesV1Create(ctx context.Context, request api.MakeupClassesV1CreateRequestObject) (api.MakeupClassesV1CreateResponseObject, error) {
	domainMC := toDomainMakeupClassFromRequest(*request.Body)
	created, err := h.makeupClassSvc.Create(ctx, domainMC)
	if err != nil {
		return nil, err
	}
	res, err := makeupClassToAPI(created)
	if err != nil {
		return nil, err
	}
	return api.MakeupClassesV1Create201JSONResponse{MakeupClass: res}, nil
}
