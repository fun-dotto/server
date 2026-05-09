package domain

import "time"

type PersonalCalendarItemStatus string

const (
	PersonalCalendarItemStatusNormal      PersonalCalendarItemStatus = "Normal"
	PersonalCalendarItemStatusCancelled   PersonalCalendarItemStatus = "Cancelled"
	PersonalCalendarItemStatusMakeup      PersonalCalendarItemStatus = "Makeup"
	PersonalCalendarItemStatusRoomChanged PersonalCalendarItemStatus = "RoomChanged"
)

type PersonalCalendarItem struct {
	Date    time.Time
	Period  Period
	Subject Subject
	Rooms   []Room
	Status  PersonalCalendarItemStatus
}
