package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/database"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimetableItemRepository struct {
	db *gorm.DB
}

func NewTimetableItemRepository(db *gorm.DB) *TimetableItemRepository {
	return &TimetableItemRepository{db: db}
}

func (r *TimetableItemRepository) timetableItemPreload(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Subject.Faculties.Faculty").
		Preload("Subject.EligibleAttributes").
		Preload("Subject.Requirements").
		Preload("Rooms.Room")
}

func (r *TimetableItemRepository) List(ctx context.Context, filter domain.TimetableItemListFilter) ([]domain.TimetableItem, error) {
	query := r.timetableItemPreload(r.db.WithContext(ctx)).
		Joins("JOIN subjects ON subjects.id = timetable_items.subject_id")

	if filter.Year != nil {
		query = query.Where("subjects.year = ?", *filter.Year)
	}
	if len(filter.Semesters) > 0 {
		semesters := make([]string, len(filter.Semesters))
		for i, s := range filter.Semesters {
			semesters[i] = string(s)
		}
		query = query.Where("subjects.semester IN ?", semesters)
	}

	query = query.Order(
		"CASE timetable_items.day_of_week " +
			"WHEN 'Monday' THEN 1 " +
			"WHEN 'Tuesday' THEN 2 " +
			"WHEN 'Wednesday' THEN 3 " +
			"WHEN 'Thursday' THEN 4 " +
			"WHEN 'Friday' THEN 5 " +
			"WHEN 'Saturday' THEN 6 " +
			"WHEN 'Sunday' THEN 7 " +
			"ELSE 8 END ASC").
		Order("timetable_items.period ASC").
		Order("subjects.syllabus_id ASC")

	var records []database.TimetableItem
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	items := make([]domain.TimetableItem, len(records))
	for i, rec := range records {
		items[i] = database.TimetableItemToDomain(rec)
	}
	return items, nil
}

func (r *TimetableItemRepository) Create(ctx context.Context, item domain.TimetableItem) (domain.TimetableItem, error) {
	record := database.TimetableItemFromDomain(item)
	record.ID = uuid.New().String()
	for i := range record.Rooms {
		record.Rooms[i].ID = uuid.New().String()
		record.Rooms[i].TimetableItemID = record.ID
	}

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.TimetableItem{}, err
	}

	var created database.TimetableItem
	if err := r.timetableItemPreload(r.db.WithContext(ctx)).First(&created, "id = ?", record.ID).Error; err != nil {
		return domain.TimetableItem{}, err
	}
	return database.TimetableItemToDomain(created), nil
}

func (r *TimetableItemRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record database.TimetableItem
		if err := tx.First(&record, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Where("timetable_item_id = ?", id).Delete(&database.TimetableItemRoom{}).Error; err != nil {
			return err
		}

		result := tx.Delete(&record)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
