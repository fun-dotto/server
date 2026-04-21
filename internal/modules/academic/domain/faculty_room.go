package domain

type FacultyRoom struct {
	ID      string
	Faculty Faculty
	Room    Room
	Year    int
}

type FacultyRoomListFilter struct {
	Year *int
}
