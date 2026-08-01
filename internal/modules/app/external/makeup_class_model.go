package external

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDomainMakeupClass は外部APIのMakeupClassをDomainに変換する
func ToDomainMakeupClass(m academic_api.MakeupClass) domain.MakeupClass {
	return domain.MakeupClass{
		ID:      m.Id,
		Comment: m.Comment,
		Date:    m.Date.Time,
		Period:  domain.Period(m.Period),
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalMakeupClassQuery はDomainのMakeupClassQueryを外部APIのパラメータに変換する
func ToExternalMakeupClassQuery(q domain.MakeupClassQuery) *academic_api.MakeupClassesV1ListParams {
	params := &academic_api.MakeupClassesV1ListParams{
		SubjectIds: q.SubjectIDs,
	}
	if q.From != nil {
		d := openapi_types.Date{Time: *q.From}
		params.From = &d
	}
	if q.Until != nil {
		d := openapi_types.Date{Time: *q.Until}
		params.Until = &d
	}
	return params
}
