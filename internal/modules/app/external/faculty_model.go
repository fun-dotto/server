package external

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// ToDomainFaculty は外部APIのFacultyをDomainのFacultyに変換する
func ToDomainFaculty(f academic_api.Faculty) domain.Faculty {
	return domain.Faculty{
		ID:    f.Id,
		Name:  f.Name,
		Email: f.Email,
	}
}

// ToDomainFaculties は外部APIのFaculty一覧をDomainのFaculty一覧に変換する
func ToDomainFaculties(faculties []academic_api.Faculty) []domain.Faculty {
	result := make([]domain.Faculty, len(faculties))
	for i, f := range faculties {
		result[i] = ToDomainFaculty(f)
	}
	return result
}
