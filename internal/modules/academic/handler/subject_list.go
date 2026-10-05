package handler

import (
	"context"

	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func (h *Handler) SubjectsV1List(ctx context.Context, request api.SubjectsV1ListRequestObject) (api.SubjectsV1ListResponseObject, error) {
	filter := buildSubjectListFilter(request.Params)

	var sort *domain.SubjectListSort
	if request.Params.UserId != nil {
		user, found, err := h.userSvc.FindByID(ctx, *request.Params.UserId)
		if err != nil {
			return nil, err
		}
		if found {
			sort = &domain.SubjectListSort{
				UserCourse: user.Course,
				UserGrade:  user.Grade,
			}
		}
	}

	subjects, err := h.subjectSvc.List(ctx, filter, sort)
	if err != nil {
		return nil, err
	}
	return api.SubjectsV1List200JSONResponse{Subjects: subjectsToAPI(subjects)}, nil
}
