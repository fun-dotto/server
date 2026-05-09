package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
)

func (h *Handler) CourseRegistrationsV1Create(ctx context.Context, request api.CourseRegistrationsV1CreateRequestObject) (api.CourseRegistrationsV1CreateResponseObject, error) {
	domainCR := toDomainCourseRegistrationFromRequest(*request.Body)
	created, err := h.courseRegistrationSvc.Create(ctx, domainCR)
	if err != nil {
		return nil, err
	}
	return api.CourseRegistrationsV1Create201JSONResponse{
		CourseRegistration: courseRegistrationToAPI(created),
	}, nil
}
