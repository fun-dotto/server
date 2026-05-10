package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) List(ctx context.Context, filter domain.RoomListFilter) ([]domain.Room, error) {
	var records []model.Room
	query := r.db.WithContext(ctx)
	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", parseUUIDs(filter.IDs))
	}
	if len(filter.Floors) > 0 {
		floors := make([]string, len(filter.Floors))
		for i, f := range filter.Floors {
			floors[i] = string(f)
		}
		query = query.Where("floor IN ?", floors)
	}
	if err := query.Order("floor ASC, name ASC").Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.Room, len(records))
	for i, rec := range records {
		results[i] = roomToDomain(rec)
	}
	return results, nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id string) (domain.Room, error) {
	var record model.Room
	if err := r.db.WithContext(ctx).First(&record, "id = ?", parseUUIDOrNil(id)).Error; err != nil {
		return domain.Room{}, err
	}
	return roomToDomain(record), nil
}

func (r *RoomRepository) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	record := roomFromDomain(room)
	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return domain.Room{}, err
	}
	return roomToDomain(record), nil
}

func (r *RoomRepository) Update(ctx context.Context, room domain.Room) (domain.Room, error) {
	id := parseUUIDOrNil(room.ID)
	if err := r.db.WithContext(ctx).Model(&model.Room{}).Where("id = ?", id).Updates(map[string]any{
		"name":  room.Name,
		"floor": string(room.Floor),
	}).Error; err != nil {
		return domain.Room{}, err
	}
	return r.GetByID(ctx, room.ID)
}

// ErrRoomInUse は、Room に紐づく FacultyRoom が存在し削除できないことを示す。
var ErrRoomInUse = errors.New("room is in use by faculty room assignments")

func (r *RoomRepository) Delete(ctx context.Context, id string) error {
	uid := parseUUIDOrNil(id)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var room model.Room
		if err := tx.Where("id = ?", uid).First(&room).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		var facultyRoomCount int64
		if err := tx.Model(&model.FacultyRoom{}).Where("room_id = ?", uid).Count(&facultyRoomCount).Error; err != nil {
			return err
		}
		if facultyRoomCount > 0 {
			return ErrRoomInUse
		}

		result := tx.Where("id = ?", uid).Delete(&model.Room{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
