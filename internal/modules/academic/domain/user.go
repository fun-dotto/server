package domain

type User struct {
	ID     string
	Course *CourseType
	Grade  *Grade
}
