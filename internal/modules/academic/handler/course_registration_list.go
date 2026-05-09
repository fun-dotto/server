package handler

import (
	"context"

	api "github.com/fun-dotto/academic-api/generated"
)

func (h *Handler) CourseRegistrationsV1List(ctx context.Context, request api.CourseRegistrationsV1ListRequestObject) (api.CourseRegistrationsV1ListResponseObject, error) {
	filter := buildCourseRegistrationListFilter(request.Params)
	items, err := h.courseRegistrationSvc.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return api.CourseRegistrationsV1List200JSONResponse{
		CourseRegistrations: courseRegistrationsToAPI(items),
	}, nil
}
