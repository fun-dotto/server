package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/server/internal/modules/academic/database"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ErrFacultyRoomAlreadyExists は、同一年度で同じ教員または同じ教室の教員室が
// 既に登録されている場合に返される。
var ErrFacultyRoomAlreadyExists = errors.New("faculty room already exists for the same faculty or room in the year")

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

	var records []database.FacultyRoom
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	items := make([]domain.FacultyRoom, len(records))
	for i, rec := range records {
		items[i] = database.FacultyRoomToDomain(rec)
	}
	return items, nil
}

func (r *FacultyRoomRepository) Create(ctx context.Context, fr domain.FacultyRoom) (domain.FacultyRoom, error) {
	record := database.FacultyRoomFromDomain(fr)
	if record.ID == "" {
		record.ID = uuid.New().String()
	}

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "idx_faculty_rooms_faculty_year", "idx_faculty_rooms_room_year":
				return domain.FacultyRoom{}, ErrFacultyRoomAlreadyExists
			}
		}
		return domain.FacultyRoom{}, err
	}

	var created database.FacultyRoom
	if err := r.facultyRoomPreload(r.db.WithContext(ctx)).First(&created, "id = ?", record.ID).Error; err != nil {
		return domain.FacultyRoom{}, err
	}
	return database.FacultyRoomToDomain(created), nil
}

func (r *FacultyRoomRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var record database.FacultyRoom
		if err := tx.Where("id = ?", id).First(&record).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", id).Delete(&database.FacultyRoom{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
