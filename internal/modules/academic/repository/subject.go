package repository

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) subjectPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Faculties").
		Preload("EligibleAttributes").
		Preload("Requirements")
}

func (r *SubjectRepository) List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error) {
	ctxDB := r.db.WithContext(ctx)
	query := r.subjectPreload(ctxDB)

	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", filter.IDs)
	}
	if filter.Q != nil {
		// TODO: filter.Q に含まれる LIKE ワイルドカード文字（%, _）をエスケープする。現状ユーザー入力がそのまま LIKE パターンに埋め込まれる。
		query = query.Where("name ILIKE ?", "%"+*filter.Q+"%")
	}
	if filter.Year != nil {
		query = query.Where("year = ?", *filter.Year)
	}
	if len(filter.Semester) > 0 {
		query = query.Where("semester IN ?", filter.Semester)
	}
	if len(filter.Grade) > 0 {
		query = query.Where("id IN (?)",
			ctxDB.Model(&database.SubjectEligibleAttribute{}).Select("subject_id").Where("grade IN ?", filter.Grade),
		)
	}
	if len(filter.Class) > 0 {
		query = query.Where("id IN (?)",
			ctxDB.Model(&database.SubjectEligibleAttribute{}).Select("subject_id").Where("class IN ?", filter.Class),
		)
	}
	if len(filter.Courses) > 0 {
		query = query.Where("id IN (?)",
			ctxDB.Model(&database.SubjectRequirement{}).Select("subject_id").Where("course IN ?", filter.Courses),
		)
	}
	if len(filter.RequirementType) > 0 {
		query = query.Where("id IN (?)",
			ctxDB.Model(&database.SubjectRequirement{}).Select("subject_id").Where("requirement_type IN ?", filter.RequirementType),
		)
	}
	if len(filter.Classification) > 0 {
		query = query.Where("classification IN ?", filter.Classification)
	}
	if len(filter.CulturalSubjectCategory) > 0 {
		query = query.Where("cultural_subject_category IN ?", filter.CulturalSubjectCategory)
	}

	var records []database.Subject
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.Subject, len(records))
	for i, rec := range records {
		results[i] = database.SubjectToDomain(rec)
	}
	return results, nil
}

func (r *SubjectRepository) GetByID(ctx context.Context, id string) (domain.Subject, error) {
	var record database.Subject
	if err := r.subjectPreload(r.db.WithContext(ctx)).First(&record, "id = ?", id).Error; err != nil {
		return domain.Subject{}, err
	}
	return database.SubjectToDomain(record), nil
}

func (r *SubjectRepository) GetBySyllabusID(ctx context.Context, syllabusID string) (domain.Subject, error) {
	var record database.Subject
	if err := r.subjectPreload(r.db.WithContext(ctx)).First(&record, "syllabus_id = ?", syllabusID).Error; err != nil {
		return domain.Subject{}, err
	}
	return database.SubjectToDomain(record), nil
}

func (r *SubjectRepository) Upsert(ctx context.Context, subject domain.Subject) (domain.Subject, error) {
	record := database.SubjectFromDomain(subject)
	record.ID = uuid.New().String()

	for i := range record.Faculties {
		record.Faculties[i].ID = uuid.New().String()
		record.Faculties[i].SubjectID = record.ID
	}
	for i := range record.EligibleAttributes {
		record.EligibleAttributes[i].ID = uuid.New().String()
		record.EligibleAttributes[i].SubjectID = record.ID
	}
	for i := range record.Requirements {
		record.Requirements[i].ID = uuid.New().String()
		record.Requirements[i].SubjectID = record.ID
	}

	// TODO: Upsert のたびに子テーブル（Faculties, EligibleAttributes, Requirements）を全件 DELETE → INSERT している。
	// データ量が増えた場合のパフォーマンスに注意。差分更新の検討が必要。
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// INSERT ... ON CONFLICT (syllabus_id) DO UPDATE で原子的に upsert
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "syllabus_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "year", "semester", "credit", "classification", "cultural_subject_category", "updated_at",
			}),
		}).Omit("Faculties", "EligibleAttributes", "Requirements").Create(&record).Error; err != nil {
			return err
		}

		// ON CONFLICT で UPDATE された場合、record.ID は新規生成した値のままなので
		// 実際の ID を取得し直す
		var actual database.Subject
		if err := tx.Select("id").Where("syllabus_id = ?", record.SyllabusID).First(&actual).Error; err != nil {
			return err
		}
		record.ID = actual.ID

		// 子テーブルの SubjectID を実際の ID に合わせる
		for i := range record.Faculties {
			record.Faculties[i].SubjectID = record.ID
		}
		for i := range record.EligibleAttributes {
			record.EligibleAttributes[i].SubjectID = record.ID
		}
		for i := range record.Requirements {
			record.Requirements[i].SubjectID = record.ID
		}

		// 1:N 子テーブルを差し替え
		if err := tx.Where("subject_id = ?", record.ID).Delete(&database.SubjectFaculty{}).Error; err != nil {
			return err
		}
		if len(record.Faculties) > 0 {
			if err := tx.Create(&record.Faculties).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("subject_id = ?", record.ID).Delete(&database.SubjectEligibleAttribute{}).Error; err != nil {
			return err
		}
		if len(record.EligibleAttributes) > 0 {
			if err := tx.Create(&record.EligibleAttributes).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("subject_id = ?", record.ID).Delete(&database.SubjectRequirement{}).Error; err != nil {
			return err
		}
		if len(record.Requirements) > 0 {
			if err := tx.Create(&record.Requirements).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.Subject{}, err
	}

	return r.GetByID(ctx, record.ID)
}

func (r *SubjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record database.Subject
		if err := tx.First(&record, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Where("subject_id = ?", id).Delete(&database.SubjectFaculty{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subject_id = ?", id).Delete(&database.SubjectEligibleAttribute{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subject_id = ?", id).Delete(&database.SubjectRequirement{}).Error; err != nil {
			return err
		}

		return tx.Delete(&record).Error
	})
}
