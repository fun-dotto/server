package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

// ListSyllabusNames は科目名照合・subject_id 解決に使う (id, syllabus_id, name) を全件返す。
func (r *SubjectRepository) ListSyllabusNames(ctx context.Context) ([]domain.SubjectRef, error) {
	var records []model.Subject
	if err := r.db.WithContext(ctx).Select("id", "syllabus_id", "name").Find(&records).Error; err != nil {
		return nil, err
	}

	refs := make([]domain.SubjectRef, len(records))
	for i, rec := range records {
		refs[i] = domain.SubjectRef{
			ID:         rec.ID.String(),
			SyllabusID: rec.SyllabusID,
			Name:       rec.Name,
		}
	}
	return refs, nil
}
