package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/shared/model"
)

func (r *CourseRegistrationRepository) ListUserIDsBySubject(ctx context.Context, subjectID string) ([]string, error) {
	var userIDs []string
	if err := r.db.WithContext(ctx).
		Model(&model.CourseRegistration{}).
		Where("subject_id = ?", subjectID).
		Distinct("user_id").
		Pluck("user_id", &userIDs).Error; err != nil {
		return nil, err
	}
	return userIDs, nil
}
