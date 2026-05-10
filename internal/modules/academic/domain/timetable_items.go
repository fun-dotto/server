package domain

type DayOfWeek string

const (
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
)

type Period string

const (
	PeriodPeriod1 Period = "Period1"
	PeriodPeriod2 Period = "Period2"
	PeriodPeriod3 Period = "Period3"
	PeriodPeriod4 Period = "Period4"
	PeriodPeriod5 Period = "Period5"
	PeriodPeriod6 Period = "Period6"
)

type TimetableSlot struct {
	DayOfWeek DayOfWeek
	Period    Period
}

type TimetableItem struct {
	ID      string
	Subject Subject
	Slot    *TimetableSlot
	Rooms   []Room
}

type TimetableItemListFilter struct {
	Year      *int
	Semesters []CourseSemester
}
