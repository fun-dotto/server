package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/user"
)

func (h *Handler) UsersV1Upsert(ctx context.Context, request api.UsersV1UpsertRequestObject) (api.UsersV1UpsertResponseObject, error) {
	domainUser := toDomainUser(request.Id, *request.Body)

	user, err := h.userService.UpsertUser(ctx, domainUser)
	if err != nil {
		return nil, err
	}

	return api.UsersV1Upsert200JSONResponse{
		User: toAPIUser(user),
	}, nil
}
