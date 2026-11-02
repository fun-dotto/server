package domain

type Term struct {
	Year     int
	Semester CourseSemester
	Start    string
	End      string
}

func (t Term) Contains(date string) bool {
	return t.Start <= date && date <= t.End
}
