package service

import (
	"context"
	"sort"
	"time"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

type personalCalendarItemCourseRegistrationRepository interface {
	List(ctx context.Context, filter domain.CourseRegistrationListFilter) ([]domain.CourseRegistration, error)
}

type personalCalendarItemTimetableItemRepository interface {
	List(ctx context.Context, filter domain.TimetableItemListFilter) ([]domain.TimetableItem, error)
}

type personalCalendarItemCancelledClassRepository interface {
	List(ctx context.Context, filter domain.CancelledClassListFilter) ([]domain.CancelledClass, error)
}

type personalCalendarItemMakeupClassRepository interface {
	List(ctx context.Context, filter domain.MakeupClassListFilter) ([]domain.MakeupClass, error)
}

type personalCalendarItemRoomChangeRepository interface {
	List(ctx context.Context, filter domain.RoomChangeListFilter) ([]domain.RoomChange, error)
}

type PersonalCalendarItemService struct {
	courseRegistrationRepo personalCalendarItemCourseRegistrationRepository
	timetableItemRepo      personalCalendarItemTimetableItemRepository
	cancelledClassRepo     personalCalendarItemCancelledClassRepository
	makeupClassRepo        personalCalendarItemMakeupClassRepository
	roomChangeRepo         personalCalendarItemRoomChangeRepository
	substituteDayMap       map[string]domain.DayOfWeek
	holidaySet             map[string]struct{}
}

func NewPersonalCalendarItemService(
	courseRegistrationRepo personalCalendarItemCourseRegistrationRepository,
	timetableItemRepo personalCalendarItemTimetableItemRepository,
	cancelledClassRepo personalCalendarItemCancelledClassRepository,
	makeupClassRepo personalCalendarItemMakeupClassRepository,
	roomChangeRepo personalCalendarItemRoomChangeRepository,
	substituteDayMap map[string]domain.DayOfWeek,
	holidaySet map[string]struct{},
) *PersonalCalendarItemService {
	return &PersonalCalendarItemService{
		courseRegistrationRepo: courseRegistrationRepo,
		timetableItemRepo:      timetableItemRepo,
		cancelledClassRepo:     cancelledClassRepo,
		makeupClassRepo:        makeupClassRepo,
		roomChangeRepo:         roomChangeRepo,
		substituteDayMap:       substituteDayMap,
		holidaySet:             holidaySet,
	}
}

func (s *PersonalCalendarItemService) List(
	ctx context.Context, userID string, dates []time.Time,
) ([]domain.PersonalCalendarItem, error) {
	// datesが空の場合、全DBクエリをスキップして早期リターン
	if len(dates) == 0 {
		return []domain.PersonalCalendarItem{}, nil
	}

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

	subjectIDSlice := make([]string, 0, len(registeredSubjectIDs))
	for id := range registeredSubjectIDs {
		subjectIDSlice = append(subjectIDSlice, id)
	}

	// 履修登録がない場合、SubjectIDsフィルタが効かず全件取得になるためスキップして早期リターン
	if len(subjectIDSlice) == 0 {
		return []domain.PersonalCalendarItem{}, nil
	}

	year, semesters := determineSemestersFromDates(dates)
	timetableItems, err := s.timetableItemRepo.List(ctx, domain.TimetableItemListFilter{
		Year:      year,
		Semesters: semesters,
	})
	if err != nil {
		return nil, err
	}

	dayToTimetableItems := make(map[domain.DayOfWeek][]domain.TimetableItem)
	subjectToRooms := make(map[string][]domain.Room)
	subjectToBestItem := make(map[string]domain.TimetableItem) // 補講用に最適な時間割を保持

	for _, item := range timetableItems {
		if item.Slot == nil {
			continue
		}
		if _, ok := registeredSubjectIDs[item.Subject.ID]; !ok {
			continue
		}
		dayToTimetableItems[item.Slot.DayOfWeek] = append(dayToTimetableItems[item.Slot.DayOfWeek], item)

		// 補講用に最適な時間割を選択（空でないRoomsを優先し、同点なら曜日/時限が最小のものを選ぶ）
		bestItem, exists := subjectToBestItem[item.Subject.ID]
		if !exists {
			subjectToBestItem[item.Subject.ID] = item
		} else {
			// 空でないRoomsを優先
			currentHasRooms := len(item.Rooms) > 0
			bestHasRooms := len(bestItem.Rooms) > 0

			shouldReplace := false
			if currentHasRooms && !bestHasRooms {
				shouldReplace = true
			} else if currentHasRooms == bestHasRooms {
				// 同条件なら曜日/時限が辞書順で小さい方を選ぶ
				if item.Slot.DayOfWeek < bestItem.Slot.DayOfWeek {
					shouldReplace = true
				} else if item.Slot.DayOfWeek == bestItem.Slot.DayOfWeek && item.Slot.Period < bestItem.Slot.Period {
					shouldReplace = true
				}
			}

			if shouldReplace {
				subjectToBestItem[item.Subject.ID] = item
			}
		}
	}

	// 選択された最適な時間割から教室情報を抽出（重複削除）
	for subjectID, item := range subjectToBestItem {
		seenRoomIDs := make(map[string]struct{})
		uniqueRooms := make([]domain.Room, 0, len(item.Rooms))
		for _, room := range item.Rooms {
			if _, exists := seenRoomIDs[room.ID]; !exists {
				seenRoomIDs[room.ID] = struct{}{}
				uniqueRooms = append(uniqueRooms, room)
			}
		}
		subjectToRooms[subjectID] = uniqueRooms
	}

	// datesからFrom/Untilを算出
	from, until := dateRange(dates)

	// 休講・補講・教室変更を取得
	cancelledClasses, err := s.cancelledClassRepo.List(ctx, domain.CancelledClassListFilter{
		SubjectIDs: subjectIDSlice,
		From:       from,
		Until:      until,
	})
	if err != nil {
		return nil, err
	}

	makeupClasses, err := s.makeupClassRepo.List(ctx, domain.MakeupClassListFilter{
		SubjectIDs: subjectIDSlice,
		From:       from,
		Until:      until,
	})
	if err != nil {
		return nil, err
	}

	roomChanges, err := s.roomChangeRepo.List(ctx, domain.RoomChangeListFilter{
		SubjectIDs: subjectIDSlice,
		From:       from,
		Until:      until,
	})
	if err != nil {
		return nil, err
	}

	// date+period+subjectID をキーとするmapを構築
	type calendarKey struct {
		date      string
		period    domain.Period
		subjectID string
	}

	cancelledMap := make(map[calendarKey]struct{})
	for _, cc := range cancelledClasses {
		cancelledMap[calendarKey{date: cc.Date, period: cc.Period, subjectID: cc.Subject.ID}] = struct{}{}
	}

	roomChangeMap := make(map[calendarKey]domain.RoomChange)
	for _, rc := range roomChanges {
		roomChangeMap[calendarKey{date: rc.Date, period: rc.Period, subjectID: rc.Subject.ID}] = rc
	}

	makeupMap := make(map[calendarKey]domain.MakeupClass)
	for _, mc := range makeupClasses {
		makeupMap[calendarKey{date: mc.Date, period: mc.Period, subjectID: mc.Subject.ID}] = mc
	}

	// 優先度: Cancelled > Makeup > RoomChanged > Normal
	// 同一 date/period/subject に複数の状態が存在する場合に重複アイテムが返らないよう map で管理する
	resultMap := make(map[calendarKey]domain.PersonalCalendarItem)

	// 時間割ベースのアイテム生成
	for _, date := range dates {
		dateStr := date.Format("2006-01-02")

		// 休日の場合はスキップ
		if _, isHoliday := s.holidaySet[dateStr]; isHoliday {
			continue
		}

		dow, ok := s.substituteDayMap[dateStr]
		if !ok {
			dow = weekdayToDayOfWeek(date.Weekday())
		}
		items := dayToTimetableItems[dow]
		for _, item := range items {
			key := calendarKey{date: dateStr, period: item.Slot.Period, subjectID: item.Subject.ID}

			if _, cancelled := cancelledMap[key]; cancelled {
				resultMap[key] = domain.PersonalCalendarItem{
					Date:    date,
					Period:  item.Slot.Period,
					Subject: item.Subject,
					Rooms:   item.Rooms,
					Status:  domain.PersonalCalendarItemStatusCancelled,
				}
				continue
			}

			if rc, changed := roomChangeMap[key]; changed {
				resultMap[key] = domain.PersonalCalendarItem{
					Date:    date,
					Period:  item.Slot.Period,
					Subject: item.Subject,
					Rooms:   []domain.Room{rc.NewRoom},
					Status:  domain.PersonalCalendarItemStatusRoomChanged,
				}
				continue
			}

			resultMap[key] = domain.PersonalCalendarItem{
				Date:    date,
				Period:  item.Slot.Period,
				Subject: item.Subject,
				Rooms:   item.Rooms,
				Status:  domain.PersonalCalendarItemStatusNormal,
			}
		}

		// 補講アイテムの追加
		// Cancelled が既にある同キーは上書きしない（Cancelled > Makeup の優先度）
		for key, mc := range makeupMap {
			if key.date != dateStr {
				continue
			}
			if existing, exists := resultMap[key]; exists && existing.Status == domain.PersonalCalendarItemStatusCancelled {
				continue
			}
			// 補講の教室情報はDBに保存されていないため、
			// 該当科目の通常時間割から選択された最適な教室を代替として使用する。
			// （空でないRoomsを優先し、同点なら曜日/時限が辞書順で最小のものを選択）
			// 実際の補講教室とは異なる可能性がある。
			rooms, ok := subjectToRooms[mc.Subject.ID]
			if !ok {
				rooms = []domain.Room{}
			}
			resultMap[key] = domain.PersonalCalendarItem{
				Date:    date,
				Period:  mc.Period,
				Subject: mc.Subject,
				Rooms:   rooms,
				Status:  domain.PersonalCalendarItemStatusMakeup,
			}
		}
	}

	result := make([]domain.PersonalCalendarItem, 0, len(resultMap))
	for _, item := range resultMap {
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		if !result[i].Date.Equal(result[j].Date) {
			return result[i].Date.Before(result[j].Date)
		}
		if result[i].Period != result[j].Period {
			return result[i].Period < result[j].Period
		}
		return result[i].Subject.ID < result[j].Subject.ID
	})

	return result, nil
}

func dateRange(dates []time.Time) (*time.Time, *time.Time) {
	if len(dates) == 0 {
		return nil, nil
	}
	minDate := dates[0]
	maxDate := dates[0]
	for _, d := range dates[1:] {
		if d.Before(minDate) {
			minDate = d
		}
		if d.After(maxDate) {
			maxDate = d
		}
	}
	return &minDate, &maxDate
}

func determineSemestersFromDates(dates []time.Time) (*int, []domain.CourseSemester) {
	if len(dates) == 0 {
		return nil, nil
	}

	yearMap := make(map[int]struct{})
	semesterMap := make(map[domain.CourseSemester]struct{})

	for _, date := range dates {
		yearMap[date.Year()] = struct{}{}

		month := date.Month()
		if month >= 4 && month <= 9 {
			semesterMap[domain.CourseSemesterH1] = struct{}{}
			semesterMap[domain.CourseSemesterQ1] = struct{}{}
			semesterMap[domain.CourseSemesterQ2] = struct{}{}
			semesterMap[domain.CourseSemesterAllYear] = struct{}{}
		} else {
			semesterMap[domain.CourseSemesterH2] = struct{}{}
			semesterMap[domain.CourseSemesterQ3] = struct{}{}
			semesterMap[domain.CourseSemesterQ4] = struct{}{}
			semesterMap[domain.CourseSemesterAllYear] = struct{}{}
		}
	}

	var year *int
	if len(yearMap) == 1 {
		for y := range yearMap {
			y2 := y
			year = &y2
			break
		}
	}

	semesters := make([]domain.CourseSemester, 0, len(semesterMap))
	for sem := range semesterMap {
		semesters = append(semesters, sem)
	}

	return year, semesters
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
