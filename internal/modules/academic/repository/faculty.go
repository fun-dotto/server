package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/academic-api/internal/database"
	"github.com/fun-dotto/academic-api/internal/domain"
	"gorm.io/gorm"
)

type FacultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) *FacultyRepository {
	return &FacultyRepository{db: db}
}

func (r *FacultyRepository) List(ctx context.Context, ids []string) ([]domain.Faculty, error) {
	var dbFaculties []database.Faculty
	query := r.db.WithContext(ctx)
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	if err := query.Order("email ASC").Find(&dbFaculties).Error; err != nil {
		return nil, err
	}

	domainFaculties := make([]domain.Faculty, len(dbFaculties))
	for i, dbFaculty := range dbFaculties {
		domainFaculties[i] = database.FacultyToDomain(dbFaculty)
	}

	return domainFaculties, nil
}

func (r *FacultyRepository) GetByID(ctx context.Context, id string) (domain.Faculty, error) {
	var dbFaculty database.Faculty
	if err := r.db.WithContext(ctx).First(&dbFaculty, "id = ?", id).Error; err != nil {
		return domain.Faculty{}, err
	}
	return database.FacultyToDomain(dbFaculty), nil
}

func (r *FacultyRepository) Create(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	// TODO: faculty.ID が空の場合に uuid.New().String() で採番する（空のまま渡すと primary key 制約違反で INSERT 失敗）
	dbFaculty := database.FacultyFromDomain(faculty)
	if err := r.db.WithContext(ctx).Create(&dbFaculty).Error; err != nil {
		return domain.Faculty{}, err
	}
	return database.FacultyToDomain(dbFaculty), nil
}

func (r *FacultyRepository) Update(ctx context.Context, faculty domain.Faculty) (domain.Faculty, error) {
	dbFaculty := database.FacultyFromDomain(faculty)
	if err := r.db.WithContext(ctx).Model(&database.Faculty{}).Where("id = ?", faculty.ID).Updates(map[string]interface{}{
		"name":  dbFaculty.Name,
		"email": dbFaculty.Email,
	}).Error; err != nil {
		return domain.Faculty{}, err
	}
	return r.GetByID(ctx, faculty.ID)
}

func (r *FacultyRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Ensure the faculty exists before deleting dependent rows to keep not-found deletes side-effect free.
		var faculty database.Faculty
		if err := tx.Where("id = ?", id).First(&faculty).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return err
		}

		if err := tx.Where("faculty_id = ?", id).Delete(&database.SubjectFaculty{}).Error; err != nil {
			return err
		}

		if err := tx.Where("faculty_id = ?", id).Delete(&database.FacultyRoom{}).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", id).Delete(&database.Faculty{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
