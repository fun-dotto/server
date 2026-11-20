package handler

import (
	"context"
	"errors"

	"gorm.io/gorm"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) CourseRegistrationsV1Delete(ctx context.Context, request api.CourseRegistrationsV1DeleteRequestObject) (api.CourseRegistrationsV1DeleteResponseObject, error) {
	err := h.courseRegistrationSvc.Delete(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.CourseRegistrationsV1Delete404Response{}, nil
		}
		return nil, err
	}
	return api.CourseRegistrationsV1Delete204Response{}, nil
}
