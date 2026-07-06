package domain

// FacultyRoomImportRow は faculty_rooms 取り込み CSV の 1 行（year は --faculties オプションで指定）。
type FacultyRoomImportRow struct {
	Year     int
	Email    string
	RoomName string
}

// FacultyRoomImportSummary は faculty_rooms 取り込みの結果集計。
type FacultyRoomImportSummary struct {
	Inserted        int
	SkippedBlank    int
	UnmatchedEmails []string
	UnmatchedRooms  []string
}

// FacultyRoomInsert は email / room_name を faculty_id / room_id に解決済みの
// faculty_rooms 挿入対象 1 件を表す。
type FacultyRoomInsert struct {
	FacultyID string
	RoomID    string
	Year      int
}
