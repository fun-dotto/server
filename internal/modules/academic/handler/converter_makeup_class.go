package handler

import (
	"fmt"
	"time"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func makeupClassToAPI(d domain.MakeupClass) (api.MakeupClass, error) {
	date, err := time.Parse("2006-01-02", d.Date)
	if err != nil {
		return api.MakeupClass{}, fmt.Errorf("failed to parse makeup class date %q: %w", d.Date, err)
	}
	return api.MakeupClass{
		Id:      d.ID,
		Subject: subjectToAPI(d.Subject),
		Date:    openapi_types.Date{Time: date},
		Period:  api.DottoFoundationV1Period(d.Period),
		Comment: d.Comment,
	}, nil
}

func makeupClassesToAPI(ds []domain.MakeupClass) ([]api.MakeupClass, error) {
	result := make([]api.MakeupClass, len(ds))
	for i, d := range ds {
		r, err := makeupClassToAPI(d)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func buildMakeupClassListFilter(params api.MakeupClassesV1ListParams) domain.MakeupClassListFilter {
	filter := domain.MakeupClassListFilter{}
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

func toDomainMakeupClassFromRequest(req api.MakeupClassRequest) domain.MakeupClass {
	return domain.MakeupClass{
		Subject: domain.Subject{ID: req.SubjectId},
		Date:    req.Date.Format("2006-01-02"),
		Period:  domain.Period(req.Period),
		Comment: req.Comment,
	}
}
