package handler

import (
	"context"

	api "github.com/fun-dotto/user-api/generated"
)

func (h *Handler) UsersV1List(ctx context.Context, _ api.UsersV1ListRequestObject) (api.UsersV1ListResponseObject, error) {
	users, err := h.userService.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return api.UsersV1List200JSONResponse{
		Users: toAPIUsers(users),
	}, nil
}
