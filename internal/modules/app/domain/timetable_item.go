package domain

// TimetableItem 時間割アイテム
type TimetableItem struct {
	ID      string
	Slot    *TimetableSlot
	Rooms   []Room
	Subject Subject
}

// DayOfWeek 曜日
type DayOfWeek string

const (
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
)

// Period 時限
type Period string

const (
	PeriodPeriod1 Period = "Period1"
	PeriodPeriod2 Period = "Period2"
	PeriodPeriod3 Period = "Period3"
	PeriodPeriod4 Period = "Period4"
	PeriodPeriod5 Period = "Period5"
	PeriodPeriod6 Period = "Period6"
)

// Floor 教室の階数
type Floor string

const (
	Floor1       Floor = "Floor1"
	Floor2       Floor = "Floor2"
	Floor3       Floor = "Floor3"
	Floor4       Floor = "Floor4"
	Floor5       Floor = "Floor5"
	Floor6       Floor = "Floor6"
	Floor7       Floor = "Floor7"
	FloorVirtual Floor = "Virtual"
)

// TimetableSlot 時間割の曜日・時限
type TimetableSlot struct {
	DayOfWeek DayOfWeek
	Period    Period
}

// Room 教室
type Room struct {
	ID    string
	Name  string
	Floor Floor
}

// TimetableItemQuery 時間割アイテム検索クエリ
type TimetableItemQuery struct {
	Semesters []CourseSemester
	Year      *int
}
