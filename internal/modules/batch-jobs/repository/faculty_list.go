package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

// ListFaculties は email 照合に使う (id, email) を全件返す。
func (r *FacultyRepository) ListFaculties(ctx context.Context) ([]domain.Faculty, error) {
	var records []model.Faculty
	if err := r.db.WithContext(ctx).Select("id", "email").Find(&records).Error; err != nil {
		return nil, err
	}

	faculties := make([]domain.Faculty, len(records))
	for i, rec := range records {
		faculties[i] = domain.Faculty{ID: rec.ID.String(), Email: rec.Email}
	}
	return faculties, nil
}
