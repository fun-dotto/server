package domain

type SubjectTargetClass struct {
	Grade Grade
	Class *Class // 修士・博士課程対象の場合はnil
}

type SubjectRequirement struct {
	Course          Course
	RequirementType SubjectRequirementType
}

type Subject struct {
	ID                      string
	Name                    string
	Faculty                 Faculty
	Semester                CourseSemester
	DayOfWeekTimetableSlots []DayOfWeekTimetableSlot
	EligibleAttributes      []SubjectTargetClass
	Requirements            []SubjectRequirement
	Categories              []SubjectCategory
	SyllabusID              string
}
