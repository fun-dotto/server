package domain

import "time"

// PersonalCalendarItem 個人カレンダーアイテム
type PersonalCalendarItem struct {
	Date    time.Time
	Period  Period
	Rooms   []Room
	Status  PersonalCalendarItemStatus
	Subject Subject
}

// PersonalCalendarItemStatus 個人カレンダーアイテムの状態
type PersonalCalendarItemStatus string

const (
	PersonalCalendarItemStatusNormal      PersonalCalendarItemStatus = "Normal"
	PersonalCalendarItemStatusCancelled   PersonalCalendarItemStatus = "Cancelled"
	PersonalCalendarItemStatusMakeup      PersonalCalendarItemStatus = "Makeup"
	PersonalCalendarItemStatusRoomChanged PersonalCalendarItemStatus = "RoomChanged"
)
