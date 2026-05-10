package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) subjectPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Faculties.Faculty").
		Preload("EligibleAttributes").
		Preload("Requirements")
}

func (r *SubjectRepository) List(ctx context.Context, filter domain.SubjectListFilter) ([]domain.Subject, error) {
	ctxDB := r.db.WithContext(ctx)
	query := r.subjectPreload(ctxDB)

	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", parseUUIDs(filter.IDs))
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
	attrSubQuery := ctxDB.Model(&model.SubjectEligibleAttribute{}).Select("subject_id")
	hasAttrFilter := false
	if len(filter.Grade) > 0 {
		attrSubQuery = attrSubQuery.Where("grade IN ?", filter.Grade)
		hasAttrFilter = true
	}
	if len(filter.Class) > 0 {
		attrSubQuery = attrSubQuery.Where("class IN ?", filter.Class)
		hasAttrFilter = true
	}
	if hasAttrFilter {
		query = query.Where("id IN (?)", attrSubQuery)
	}
	reqSubQuery := ctxDB.Model(&model.SubjectRequirement{}).Select("subject_id")
	hasReqFilter := false
	if len(filter.Courses) > 0 {
		reqSubQuery = reqSubQuery.Where("course IN ?", filter.Courses)
		hasReqFilter = true
	}
	if len(filter.RequirementType) > 0 {
		reqSubQuery = reqSubQuery.Where("requirement_type IN ?", filter.RequirementType)
		hasReqFilter = true
	}
	if hasReqFilter {
		query = query.Where("id IN (?)", reqSubQuery)
	}
	if len(filter.Classification) > 0 {
		query = query.Where("classification IN ?", filter.Classification)
	}
	if len(filter.CulturalSubjectCategory) > 0 {
		query = query.Where("cultural_subject_category IN ?", filter.CulturalSubjectCategory)
	}

	query = query.Order("syllabus_id ASC")

	var records []model.Subject
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	results := make([]domain.Subject, len(records))
	for i, rec := range records {
		results[i] = subjectToDomain(rec)
	}
	return results, nil
}

func (r *SubjectRepository) GetByID(ctx context.Context, id string) (domain.Subject, error) {
	var record model.Subject
	if err := r.subjectPreload(r.db.WithContext(ctx)).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		return domain.Subject{}, err
	}
	return subjectToDomain(record), nil
}

func (r *SubjectRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record model.Subject
		if err := tx.First(&record, "id = ?", uid).Error; err != nil {
			return err
		}

		if err := tx.Where("subject_id = ?", uid).Delete(&model.SubjectFaculty{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subject_id = ?", uid).Delete(&model.SubjectEligibleAttribute{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subject_id = ?", uid).Delete(&model.SubjectRequirement{}).Error; err != nil {
			return err
		}

		return tx.Delete(&record).Error
	})
}
