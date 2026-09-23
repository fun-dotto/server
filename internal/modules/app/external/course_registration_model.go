package external

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// ToDomainCourseRegistration は外部APIのCourseRegistrationをDomainのCourseRegistrationに変換する
func ToDomainCourseRegistration(m academic_api.CourseRegistration) domain.CourseRegistration {
	return domain.CourseRegistration{
		ID:      m.Id,
		Subject: ToDomainSubjectSummary(m.Subject),
	}
}
