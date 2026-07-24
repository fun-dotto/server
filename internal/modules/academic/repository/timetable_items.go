package repository

import (
	"context"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
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

	var records []model.TimetableItem
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.TimetableItem, len(records))
	for i, rec := range records {
		results[i] = timetableItemToDomain(rec)
	}
	return results, nil
}

func (r *TimetableItemRepository) Create(ctx context.Context, item domain.TimetableItem) (domain.TimetableItem, error) {
	roomIDs := make([]uuid.UUID, len(item.Rooms))
	for i, room := range item.Rooms {
		roomIDs[i] = parseUUIDOrNil(room.ID)
	}

	// many2many の Rooms を auto-upsert させると既存 Room レコードを書き換えてしまうため、
	// 親レコードのみ作成し、join テーブル (timetable_item_rooms) は別途差し込む。
	record := model.TimetableItem{
		SubjectID: parseUUIDOrNil(item.Subject.ID),
	}
	if item.Slot != nil {
		dow := string(item.Slot.DayOfWeek)
		p := string(item.Slot.Period)
		record.DayOfWeek = &dow
		record.Period = &p
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		for _, roomID := range roomIDs {
			if err := tx.Exec(
				"INSERT INTO timetable_item_rooms (timetable_item_id, room_id) VALUES (?, ?)",
				record.ID, roomID,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.TimetableItem{}, err
	}

	var created model.TimetableItem
	if err := r.timetableItemPreload(r.db.WithContext(ctx)).First(&created, "id = ?", record.ID).Error; err != nil {
		return domain.TimetableItem{}, err
	}
	return timetableItemToDomain(created), nil
}

func (r *TimetableItemRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record model.TimetableItem
		if err := tx.First(&record, "id = ?", uid).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM timetable_item_rooms WHERE timetable_item_id = ?", uid).Error; err != nil {
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
