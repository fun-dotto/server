package domain

// RoomRef は教室名照合に使う rooms テーブルの最小表現。
type RoomRef struct {
	ID   string
	Name string
}
