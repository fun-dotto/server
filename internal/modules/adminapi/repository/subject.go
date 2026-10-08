package repository

import (
	"context"
	"slices"

	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

// SubjectRepository は子テーブル（担当教員・対象学年・必修区分）を含めて Subject を扱う。
type SubjectRepository struct {
	*Repository[model.Subject]
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{
		Repository: newRepository[model.Subject](db, []string{"id"}, "Faculties", "EligibleAttributes", "Requirements"),
	}
}

func (r *SubjectRepository) Create(ctx context.Context, v *model.Subject) error {
	return translateError(r.db.WithContext(ctx).Omit("Syllabus", "Faculties.Faculty").Create(v).Error)
}

// Update は columns に子テーブルのフィールド名（Faculties など）が含まれる場合、
// その子テーブルの行をすべて v の内容で置き換える。
func (r *SubjectRepository) Update(ctx context.Context, key Key, v *model.Subject, columns []string) error {
	if len(columns) == 0 {
		_, err := r.Get(ctx, key)
		return err
	}
	return translateError(r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 子テーブルのみ更新する場合も updated_at を進め、存在確認を兼ねる。
		cols := slices.DeleteFunc(slices.Clone(columns), func(c string) bool {
			return c == "Faculties" || c == "EligibleAttributes" || c == "Requirements"
		})
		cols = append(cols, "UpdatedAt")
		if err := updateColumns[model.Subject](tx, key, v, cols); err != nil {
			return err
		}
		var current model.Subject
		if err := tx.Where(map[string]any(key)).Take(&current).Error; err != nil {
			return err
		}
		if slices.Contains(columns, "Faculties") {
			if err := replaceChildren(tx, current.ID.String(), v.Faculties, func(c *model.SubjectFaculty) { c.SubjectID = current.ID }); err != nil {
				return err
			}
		}
		if slices.Contains(columns, "EligibleAttributes") {
			if err := replaceChildren(tx, current.ID.String(), v.EligibleAttributes, func(c *model.SubjectEligibleAttribute) { c.SubjectID = current.ID }); err != nil {
				return err
			}
		}
		if slices.Contains(columns, "Requirements") {
			if err := replaceChildren(tx, current.ID.String(), v.Requirements, func(c *model.SubjectRequirement) { c.SubjectID = current.ID }); err != nil {
				return err
			}
		}
		return nil
	}))
}

func replaceChildren[C any](tx *gorm.DB, subjectID string, children []C, setParent func(*C)) error {
	if err := tx.Where("subject_id = ?", subjectID).Delete(new(C)).Error; err != nil {
		return err
	}
	if len(children) == 0 {
		return nil
	}
	for i := range children {
		setParent(&children[i])
	}
	return tx.Omit("Faculty").Create(&children).Error
}
