package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ErrFacultyRoomAlreadyExists は、同一年度で同じ (Faculty, Room) の組み合わせが
// 既に登録されている場合に返される。
var ErrFacultyRoomAlreadyExists = errors.New("faculty room already exists for the same faculty and room in the year")

type FacultyRoomRepository struct {
	db *gorm.DB
}

func NewFacultyRoomRepository(db *gorm.DB) *FacultyRoomRepository {
	return &FacultyRoomRepository{db: db}
}

func (r *FacultyRoomRepository) facultyRoomPreload(db *gorm.DB) *gorm.DB {
	return db.Preload("Faculty").Preload("Room")
}

func (r *FacultyRoomRepository) List(ctx context.Context, filter domain.FacultyRoomListFilter) ([]domain.FacultyRoom, error) {
	query := r.facultyRoomPreload(r.db.WithContext(ctx)).
		Joins("JOIN rooms ON rooms.id = faculty_rooms.room_id").
		Order("faculty_rooms.year ASC, rooms.floor ASC, rooms.name ASC")
	if filter.Year != nil {
		query = query.Where("faculty_rooms.year = ?", *filter.Year)
	}

	var records []model.FacultyRoom
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	results := make([]domain.FacultyRoom, len(records))
	for i, rec := range records {
		results[i] = facultyRoomToDomain(rec)
	}
	return results, nil
}

func (r *FacultyRoomRepository) Create(ctx context.Context, fr domain.FacultyRoom) (domain.FacultyRoom, error) {
	record := facultyRoomFromDomain(fr)

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.FacultyRoom{}, ErrFacultyRoomAlreadyExists
		}
		return domain.FacultyRoom{}, err
	}

	var created model.FacultyRoom
	if err := r.facultyRoomPreload(r.db.WithContext(ctx)).
		Where("faculty_id = ? AND room_id = ? AND year = ?", record.FacultyID, record.RoomID, record.Year).
		First(&created).Error; err != nil {
		return domain.FacultyRoom{}, err
	}
	return facultyRoomToDomain(created), nil
}

func (r *FacultyRoomRepository) Delete(ctx context.Context, id string) error {
	facultyID, roomID, year, err := decodeFacultyRoomID(id)
	if err != nil {
		return gorm.ErrRecordNotFound
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record model.FacultyRoom
		if err := tx.Where("faculty_id = ? AND room_id = ? AND year = ?", facultyID, roomID, year).
			First(&record).Error; err != nil {
			return err
		}

		result := tx.Where("faculty_id = ? AND room_id = ? AND year = ?", facultyID, roomID, year).
			Delete(&model.FacultyRoom{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
