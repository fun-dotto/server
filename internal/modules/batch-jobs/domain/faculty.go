package domain

// Faculty は faculty_rooms インポート時に email で突合するための最小表現。
type Faculty struct {
	ID    string
	Email string
}
