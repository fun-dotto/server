package handler

import (
	"context"

	api "github.com/fun-dotto/subject-api/generated"
)

func (h *Handler) SubjectCategoriesV1List(ctx context.Context, _ api.SubjectCategoriesV1ListRequestObject) (api.SubjectCategoriesV1ListResponseObject, error) {
	categories, err := h.categorySvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return api.SubjectCategoriesV1List200JSONResponse{SubjectCategories: subjectCategoriesToAPI(categories)}, nil
}

func (h *Handler) SubjectCategoriesV1Create(ctx context.Context, request api.SubjectCategoriesV1CreateRequestObject) (api.SubjectCategoriesV1CreateResponseObject, error) {
	category, err := h.categorySvc.Create(ctx, subjectCategoryRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.SubjectCategoriesV1Create201JSONResponse{SubjectCategory: subjectCategoryToAPI(category)}, nil
}

func (h *Handler) SubjectCategoriesV1Detail(ctx context.Context, request api.SubjectCategoriesV1DetailRequestObject) (api.SubjectCategoriesV1DetailResponseObject, error) {
	category, err := h.categorySvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.SubjectCategoriesV1Detail200JSONResponse{SubjectCategory: subjectCategoryToAPI(category)}, nil
}

func (h *Handler) SubjectCategoriesV1Update(ctx context.Context, request api.SubjectCategoriesV1UpdateRequestObject) (api.SubjectCategoriesV1UpdateResponseObject, error) {
	category, err := h.categorySvc.Update(ctx, request.Id, subjectCategoryRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.SubjectCategoriesV1Update200JSONResponse{SubjectCategory: subjectCategoryToAPI(category)}, nil
}

func (h *Handler) SubjectCategoriesV1Delete(ctx context.Context, request api.SubjectCategoriesV1DeleteRequestObject) (api.SubjectCategoriesV1DeleteResponseObject, error) {
	if err := h.categorySvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.SubjectCategoriesV1Delete204Response{}, nil
}
