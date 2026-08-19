package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) MakeupClassesV1List(ctx context.Context, request api.MakeupClassesV1ListRequestObject) (api.MakeupClassesV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toMakeupClassQuery(request.Params)

	classes, err := h.academicService.GetMakeupClasses(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get makeup classes: %w", err)
	}

	apiClasses := make([]api.MakeupClass, len(classes))
	for i, c := range classes {
		apiClasses[i] = toApiMakeupClass(c)
	}

	return api.MakeupClassesV1List200JSONResponse{
		MakeupClasses: apiClasses,
	}, nil
}

func toMakeupClassQuery(params api.MakeupClassesV1ListParams) domain.MakeupClassQuery {
	query := domain.MakeupClassQuery{
		SubjectIDs: params.SubjectIds,
	}
	if params.From != nil {
		t := params.From.Time
		query.From = &t
	}
	if params.Until != nil {
		t := params.Until.Time
		query.Until = &t
	}
	return query
}

func toApiMakeupClass(c domain.MakeupClass) api.MakeupClass {
	return api.MakeupClass{
		Id:      c.ID,
		Comment: c.Comment,
		Date:    openapi_types.Date{Time: c.Date},
		Period:  api.DottoFoundationV1Period(c.Period),
		Subject: toApiSubjectSummary(c.Subject),
	}
}
