package handler

import (
	api "github.com/fun-dotto/server/gen/academic"
	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func facultyToAPI(faculty domain.Faculty) api.Faculty {
	return api.Faculty{
		Id:    faculty.ID,
		Name:  faculty.Name,
		Email: faculty.Email,
	}
}

func facultiesToAPI(faculties []domain.Faculty) []api.Faculty {
	result := make([]api.Faculty, len(faculties))
	for i, faculty := range faculties {
		result[i] = facultyToAPI(faculty)
	}
	return result
}

func toDomainFacultyFromRequest(id string, req api.FacultyRequest) domain.Faculty {
	return domain.Faculty{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
}
