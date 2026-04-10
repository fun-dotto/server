package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/user-api/generated"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (h *Handler) UsersV1Detail(ctx context.Context, request api.UsersV1DetailRequestObject) (api.UsersV1DetailResponseObject, error) {
	user, err := h.userService.GetUserByID(request.Id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.UsersV1Detail404Response{}, nil
		}
		return nil, err
	}

	return api.UsersV1Detail200JSONResponse{
		User: toAPIUser(user),
	}, nil
}
