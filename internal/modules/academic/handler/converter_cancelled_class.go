package handler

import (
	"fmt"
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func cancelledClassToAPI(d domain.CancelledClass) (api.CancelledClass, error) {
	date, err := time.Parse("2006-01-02", d.Date)
	if err != nil {
		return api.CancelledClass{}, fmt.Errorf("failed to parse cancelled class date %q: %w", d.Date, err)
	}
	return api.CancelledClass{
		Id:      d.ID,
		Subject: subjectToAPI(d.Subject),
		Date:    openapi_types.Date{Time: date},
		Period:  api.DottoFoundationV1Period(d.Period),
		Comment: d.Comment,
	}, nil
}

func cancelledClassesToAPI(ds []domain.CancelledClass) ([]api.CancelledClass, error) {
	result := make([]api.CancelledClass, len(ds))
	for i, d := range ds {
		r, err := cancelledClassToAPI(d)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func buildCancelledClassListFilter(params api.CancelledClassesV1ListParams) domain.CancelledClassListFilter {
	filter := domain.CancelledClassListFilter{}
	if params.SubjectIds != nil {
		filter.SubjectIDs = *params.SubjectIds
	}
	if params.From != nil {
		filter.From = &params.From.Time
	}
	if params.Until != nil {
		filter.Until = &params.Until.Time
	}
	return filter
}

func toDomainCancelledClassFromRequest(req api.CancelledClassRequest) domain.CancelledClass {
	return domain.CancelledClass{
		Subject: domain.Subject{ID: req.SubjectId},
		Date:    req.Date.Format("2006-01-02"),
		Period:  domain.Period(req.Period),
		Comment: req.Comment,
	}
}
