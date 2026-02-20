package handler

import (
	"context"

	api "github.com/fun-dotto/subject-api/generated"
)

func (h *Handler) CoursesV1List(ctx context.Context, _ api.CoursesV1ListRequestObject) (api.CoursesV1ListResponseObject, error) {
	courses, err := h.courseSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return api.CoursesV1List200JSONResponse{Courses: coursesToAPI(courses)}, nil
}

func (h *Handler) CoursesV1Create(ctx context.Context, request api.CoursesV1CreateRequestObject) (api.CoursesV1CreateResponseObject, error) {
	course, err := h.courseSvc.Create(ctx, courseRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.CoursesV1Create201JSONResponse{Course: courseToAPI(course)}, nil
}

func (h *Handler) CoursesV1Detail(ctx context.Context, request api.CoursesV1DetailRequestObject) (api.CoursesV1DetailResponseObject, error) {
	course, err := h.courseSvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.CoursesV1Detail200JSONResponse{Course: courseToAPI(course)}, nil
}

func (h *Handler) CoursesV1Update(ctx context.Context, request api.CoursesV1UpdateRequestObject) (api.CoursesV1UpdateResponseObject, error) {
	course, err := h.courseSvc.Update(ctx, request.Id, courseRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.CoursesV1Update200JSONResponse{Course: courseToAPI(course)}, nil
}

func (h *Handler) CoursesV1Delete(ctx context.Context, request api.CoursesV1DeleteRequestObject) (api.CoursesV1DeleteResponseObject, error) {
	if err := h.courseSvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.CoursesV1Delete204Response{}, nil
}
