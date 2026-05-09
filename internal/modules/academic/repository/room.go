package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) List(ctx context.Context, filter domain.RoomListFilter) ([]domain.Room, error) {
	var dbRooms []database.Room
	query := r.db.WithContext(ctx)
	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", filter.IDs)
	}
	if len(filter.Floors) > 0 {
		floors := make([]string, len(filter.Floors))
		for i, f := range filter.Floors {
			floors[i] = string(f)
		}
		query = query.Where("floor IN ?", floors)
	}
	if err := query.Order("floor ASC, name ASC").Find(&dbRooms).Error; err != nil {
		return nil, err
	}

	domainRooms := make([]domain.Room, len(dbRooms))
	for i, dbRoom := range dbRooms {
		domainRooms[i] = database.RoomToDomain(dbRoom)
	}

	return domainRooms, nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id string) (domain.Room, error) {
	var dbRoom database.Room
	if err := r.db.WithContext(ctx).First(&dbRoom, "id = ?", id).Error; err != nil {
		return domain.Room{}, err
	}
	return database.RoomToDomain(dbRoom), nil
}

func (r *RoomRepository) Create(ctx context.Context, room domain.Room) (domain.Room, error) {
	// TODO: room.ID が空の場合に uuid.New().String() で採番する（空のまま渡すと primary key 制約違反で INSERT 失敗）
	dbRoom := database.RoomFromDomain(room)
	if err := r.db.WithContext(ctx).Create(&dbRoom).Error; err != nil {
		return domain.Room{}, err
	}
	return database.RoomToDomain(dbRoom), nil
}

func (r *RoomRepository) Update(ctx context.Context, room domain.Room) (domain.Room, error) {
	dbRoom := database.RoomFromDomain(room)
	if err := r.db.WithContext(ctx).Model(&database.Room{}).Where("id = ?", room.ID).Updates(map[string]interface{}{
		"name":  dbRoom.Name,
		"floor": dbRoom.Floor,
	}).Error; err != nil {
		return domain.Room{}, err
	}
	return r.GetByID(ctx, room.ID)
}

// ErrRoomInUse は、Room に紐づく FacultyRoom が存在し削除できないことを示す。
var ErrRoomInUse = errors.New("room is in use by faculty room assignments")

func (r *RoomRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var room database.Room
		if err := tx.Where("id = ?", id).First(&room).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		var facultyRoomCount int64
		if err := tx.Model(&database.FacultyRoom{}).Where("room_id = ?", id).Count(&facultyRoomCount).Error; err != nil {
			return err
		}
		if facultyRoomCount > 0 {
			return ErrRoomInUse
		}

		result := tx.Where("id = ?", id).Delete(&database.Room{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
