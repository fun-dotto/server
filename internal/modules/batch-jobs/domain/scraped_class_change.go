package domain

import "time"

// CancelledClassKind は休講コメントの接頭辞から判定した補講予定の種別を表す。
type CancelledClassKind string

const (
	CancelledClassKindMakeupScheduled  CancelledClassKind = "補講あり"
	CancelledClassKindMakeupNotPlanned CancelledClassKind = "補講なし"
	CancelledClassKindMakeupUndecided  CancelledClassKind = "補講未定"
	CancelledClassKindOther            CancelledClassKind = "その他"
)

// ScrapedCancelledClass は休講テーブルの 1 行をパースした中間表現。
type ScrapedCancelledClass struct {
	Date       time.Time
	Period     string
	LessonName string
	Campus     string
	Staff      string
	Comment    string
	Kind       CancelledClassKind
}

// ScrapedMakeupClass は補講テーブルの 1 行をパースした中間表現。
type ScrapedMakeupClass struct {
	Date       time.Time
	Period     string
	LessonName string
	Campus     string
	Staff      string
	Comment    string
	RoomName   string
}

// ScrapedRoomChange は部屋変更テーブルの 1 行をパースした中間表現。
type ScrapedRoomChange struct {
	Date         time.Time
	Period       string
	LessonName   string
	Campus       string
	Staff        string
	RoomFromName string
	RoomToName   string
}
