package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/api-template/generated"
)

func (h *Handler) FCMTokenV1Upsert(_ context.Context, _ api.FCMTokenV1UpsertRequestObject) (api.FCMTokenV1UpsertResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
