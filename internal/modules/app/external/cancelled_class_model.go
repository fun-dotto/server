package external

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ToDomainCancelledClass は外部APIのCancelledClassをDomainに変換する
func ToDomainCancelledClass(m academic_api.CancelledClass) domain.CancelledClass {
	return domain.CancelledClass{
		ID:      m.Id,
		Comment: m.Comment,
		Date:    m.Date.Time,
		Period:  domain.Period(m.Period),
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}

// ToExternalCancelledClassQuery はDomainのCancelledClassQueryを外部APIのパラメータに変換する
func ToExternalCancelledClassQuery(q domain.CancelledClassQuery) *academic_api.CancelledClassesV1ListParams {
	params := &academic_api.CancelledClassesV1ListParams{
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
