package domain

type CourseRegistration struct {
	ID      string
	UserID  string
	Subject Subject
}

type CourseRegistrationListFilter struct {
	UserID    string
	Year      *int
	Semesters []CourseSemester
}
