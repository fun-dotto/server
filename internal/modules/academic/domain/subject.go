package domain

type SubjectFaculty struct {
	Faculty   Faculty
	IsPrimary bool
}

type SubjectTargetClass struct {
	Grade Grade
	Class *Class // 修士・博士課程対象の場合はnil
}

type SubjectRequirement struct {
	Course          CourseType
	RequirementType SubjectRequirementType
}

type Subject struct {
	ID                      string
	Name                    string
	Faculties               []SubjectFaculty
	Year                    int
	Semester                CourseSemester
	Credit                  int
	Classification          SubjectClassification
	CulturalSubjectCategory CulturalSubjectCategory
	EligibleAttributes      []SubjectTargetClass
	Requirements            []SubjectRequirement
	SyllabusID              string
}

type SubjectListFilter struct {
	IDs                     []string
	Q                       *string
	Grade                   []Grade
	Courses                 []CourseType
	Class                   []Class
	Classification          []SubjectClassification
	Year                    *int
	Semester                []CourseSemester
	RequirementType         []SubjectRequirementType
	CulturalSubjectCategory []CulturalSubjectCategory

	SortByUserAttribute bool
	SortCourse          *CourseType
}
