package domain

import "time"

// Reservation 教室の予約
type Reservation struct {
	ID      string
	Room    Room
	StartAt time.Time
	EndAt   time.Time
	Title   string
}

// ReservationQuery 予約検索クエリ
type ReservationQuery struct {
	RoomIDs []string
	From    *time.Time
	Until   *time.Time
}
