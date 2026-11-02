package handler

import (
	"context"
	"fmt"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (h *Handler) CancelledClassesV1List(ctx context.Context, request api.CancelledClassesV1ListRequestObject) (api.CancelledClassesV1ListResponseObject, error) {
	if h.academicService == nil {
		return nil, errAcademicServiceNotConfigured
	}

	query := toCancelledClassQuery(request.Params)

	classes, err := h.academicService.GetCancelledClasses(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cancelled classes: %w", err)
	}

	apiClasses := make([]api.CancelledClass, len(classes))
	for i, c := range classes {
		apiClasses[i] = toApiCancelledClass(c)
	}

	return api.CancelledClassesV1List200JSONResponse{
		CancelledClasses: apiClasses,
	}, nil
}

func toCancelledClassQuery(params api.CancelledClassesV1ListParams) domain.CancelledClassQuery {
	query := domain.CancelledClassQuery{
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

func toApiCancelledClass(c domain.CancelledClass) api.CancelledClass {
	return api.CancelledClass{
		Id:      c.ID,
		Comment: c.Comment,
		Date:    openapi_types.Date{Time: c.Date},
		Period:  api.DottoFoundationV1Period(c.Period),
		Subject: toApiSubjectSummary(c.Subject),
	}
}
