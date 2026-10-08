package handler

import (
	"context"
	"log"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

// buildSubjectListSort は X-Flags で subject_sort_by_user_attributes=true が指定され、
// かつ userId のユーザーが存在する場合のみソート条件を返す。それ以外は nil（現行順）。
func (h *Handler) buildSubjectListSort(ctx context.Context, userID *string, xFlags *string) (*domain.SubjectListSort, error) {
	if userID == nil || !featureFlagEnabled(xFlags, flagSubjectSortByUserAttributes) {
		return nil, nil
	}
	user, found, err := h.userSvc.FindByID(ctx, *userID)
	if err != nil {
		return nil, err
	}
	if !found {
		// treatment 群なのに現行順になるケースを KPI 分析で切り分けられるようにログに残す。
		log.Printf("subject sort skipped: user not found: userId=%s", *userID)
		return nil, nil
	}
	return &domain.SubjectListSort{
		UserCourse: user.Course,
		UserGrade:  user.Grade,
	}, nil
}
