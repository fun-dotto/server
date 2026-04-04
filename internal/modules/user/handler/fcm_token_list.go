package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) FCMTokenV1List(_ context.Context, _ api.FCMTokenV1ListRequestObject) (api.FCMTokenV1ListResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
