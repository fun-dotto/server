package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) UsersV1List(_ context.Context, _ api.UsersV1ListRequestObject) (api.UsersV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
