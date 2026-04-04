package service

import (
	"context"
	"time"

	"github.com/fun-dotto/academic-api/internal/domain"
)

type personalCalendarItemCourseRegistrationRepository interface {
	List(ctx context.Context, filter domain.CourseRegistrationListFilter) ([]domain.CourseRegistration, error)
}

type personalCalendarItemTimetableItemRepository interface {
	List(ctx context.Context, filter domain.TimetableItemListFilter) ([]domain.TimetableItem, error)
}

type PersonalCalendarItemService struct {
	courseRegistrationRepo personalCalendarItemCourseRegistrationRepository
	timetableItemRepo     personalCalendarItemTimetableItemRepository
	substituteDayMap      map[string]domain.DayOfWeek
}

func NewPersonalCalendarItemService(
	courseRegistrationRepo personalCalendarItemCourseRegistrationRepository,
	timetableItemRepo personalCalendarItemTimetableItemRepository,
	substituteDayMap map[string]domain.DayOfWeek,
) *PersonalCalendarItemService {
	return &PersonalCalendarItemService{
		courseRegistrationRepo: courseRegistrationRepo,
		timetableItemRepo:     timetableItemRepo,
		substituteDayMap:      substituteDayMap,
	}
}

func (s *PersonalCalendarItemService) List(
	ctx context.Context, userID string, dates []time.Time,
) ([]domain.PersonalCalendarItem, error) {
	registrations, err := s.courseRegistrationRepo.List(ctx, domain.CourseRegistrationListFilter{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	registeredSubjectIDs := make(map[string]struct{}, len(registrations))
	for _, r := range registrations {
		registeredSubjectIDs[r.Subject.ID] = struct{}{}
	}

	// TODO: datesからYearとSemestersを判定してフィルタに設定する
	timetableItems, err := s.timetableItemRepo.List(ctx, domain.TimetableItemListFilter{})
	if err != nil {
		return nil, err
	}

	dayToTimetableItems := make(map[domain.DayOfWeek][]domain.TimetableItem)
	for _, item := range timetableItems {
		if item.Slot == nil {
			continue
		}
		if _, ok := registeredSubjectIDs[item.Subject.ID]; !ok {
			continue
		}
		dayToTimetableItems[item.Slot.DayOfWeek] = append(dayToTimetableItems[item.Slot.DayOfWeek], item)
	}

	var result []domain.PersonalCalendarItem
	for _, date := range dates {
		dow, ok := s.substituteDayMap[date.Format("2006-01-02")]
		if !ok {
			dow = weekdayToDayOfWeek(date.Weekday())
		}
		items := dayToTimetableItems[dow]
		for _, item := range items {
			result = append(result, domain.PersonalCalendarItem{
				Date:    date,
				Period:  item.Slot.Period,
				Subject: item.Subject,
				Rooms:   item.Rooms,
				Status:  domain.PersonalCalendarItemStatusNormal,
			})
		}
	}

	return result, nil
}

func weekdayToDayOfWeek(w time.Weekday) domain.DayOfWeek {
	switch w {
	case time.Sunday:
		return domain.DayOfWeekSunday
	case time.Monday:
		return domain.DayOfWeekMonday
	case time.Tuesday:
		return domain.DayOfWeekTuesday
	case time.Wednesday:
		return domain.DayOfWeekWednesday
	case time.Thursday:
		return domain.DayOfWeekThursday
	case time.Friday:
		return domain.DayOfWeekFriday
	case time.Saturday:
		return domain.DayOfWeekSaturday
	default:
		return domain.DayOfWeekSunday
	}
}
