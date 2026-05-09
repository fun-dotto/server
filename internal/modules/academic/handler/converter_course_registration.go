package handler

import (
	"time"

	api "github.com/fun-dotto/academic-api/generated"
	"github.com/fun-dotto/academic-api/internal/domain"
)

func courseRegistrationToAPI(d domain.CourseRegistration) api.CourseRegistration {
	return api.CourseRegistration{
		Id:      d.ID,
		UserId:  d.UserID,
		Subject: subjectToAPI(d.Subject),
	}
}

func courseRegistrationsToAPI(ds []domain.CourseRegistration) []api.CourseRegistration {
	result := make([]api.CourseRegistration, len(ds))
	for i, d := range ds {
		result[i] = courseRegistrationToAPI(d)
	}
	return result
}

func buildCourseRegistrationListFilter(params api.CourseRegistrationsV1ListParams) domain.CourseRegistrationListFilter {
	filter := domain.CourseRegistrationListFilter{
		UserID: params.UserId,
	}
	if params.Year != nil {
		filter.Year = params.Year
	} else {
		currentYear := time.Now().Year()
		filter.Year = &currentYear
	}
	semesters := make([]domain.CourseSemester, len(params.Semesters))
	for i, s := range params.Semesters {
		semesters[i] = domain.CourseSemester(s)
	}
	filter.Semesters = semesters
	return filter
}

func toDomainCourseRegistrationFromRequest(req api.CourseRegistrationRequest) domain.CourseRegistration {
	return domain.CourseRegistration{
		UserID:  req.UserId,
		Subject: domain.Subject{ID: req.SubjectId},
	}
}
