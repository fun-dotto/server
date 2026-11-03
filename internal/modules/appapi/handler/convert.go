package handler

import (
	"time"

	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoName(s string) protoreflect.Name { return protoreflect.Name(s) }

// toProtoDate は time.Time の日付部分を google.type.Date にする。
func toProtoDate(t time.Time) *date.Date {
	y, m, d := t.Date()
	return &date.Date{Year: int32(y), Month: int32(m), Day: int32(d)}
}

// fromProtoDate は google.type.Date を time.Time（UTC の 0 時）にする。不正な日付の場合は false。
func fromProtoDate(d *date.Date) (time.Time, bool) {
	if d == nil || d.GetYear() <= 0 || d.GetMonth() < 1 || d.GetMonth() > 12 || d.GetDay() < 1 {
		return time.Time{}, false
	}
	t := time.Date(int(d.GetYear()), time.Month(d.GetMonth()), int(d.GetDay()), 0, 0, 0, 0, time.UTC)
	if t.Day() != int(d.GetDay()) {
		return time.Time{}, false
	}
	return t, true
}

func toTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// bimap は domain の文字列 enum と proto の enum の対応。
type bimap[D comparable, P comparable] struct {
	toProto   map[D]P
	fromProto map[P]D
}

func newBimap[D comparable, P comparable](pairs map[D]P) bimap[D, P] {
	b := bimap[D, P]{toProto: pairs, fromProto: make(map[P]D, len(pairs))}
	for d, p := range pairs {
		b.fromProto[p] = d
	}
	return b
}

var (
	gradeMap = newBimap(map[domain.Grade]appv1.Grade{
		domain.GradeB1: appv1.Grade_GRADE_B1,
		domain.GradeB2: appv1.Grade_GRADE_B2,
		domain.GradeB3: appv1.Grade_GRADE_B3,
		domain.GradeB4: appv1.Grade_GRADE_B4,
		domain.GradeM1: appv1.Grade_GRADE_M1,
		domain.GradeM2: appv1.Grade_GRADE_M2,
		domain.GradeD1: appv1.Grade_GRADE_D1,
		domain.GradeD2: appv1.Grade_GRADE_D2,
		domain.GradeD3: appv1.Grade_GRADE_D3,
	})
	courseMap = newBimap(map[domain.Course]appv1.Course{
		domain.CourseInformationSystem: appv1.Course_COURSE_INFORMATION_SYSTEM,
		domain.CourseInformationDesign: appv1.Course_COURSE_INFORMATION_DESIGN,
		domain.CourseAdvancedICT:       appv1.Course_COURSE_ADVANCED_ICT,
		domain.CourseComplexSystem:     appv1.Course_COURSE_COMPLEX_SYSTEM,
		domain.CourseIntelligentSystem: appv1.Course_COURSE_INTELLIGENT_SYSTEM,
	})
	classMap = newBimap(map[domain.Class]appv1.Class{
		domain.ClassA: appv1.Class_CLASS_A,
		domain.ClassB: appv1.Class_CLASS_B,
		domain.ClassC: appv1.Class_CLASS_C,
		domain.ClassD: appv1.Class_CLASS_D,
		domain.ClassE: appv1.Class_CLASS_E,
		domain.ClassF: appv1.Class_CLASS_F,
		domain.ClassG: appv1.Class_CLASS_G,
		domain.ClassH: appv1.Class_CLASS_H,
		domain.ClassI: appv1.Class_CLASS_I,
		domain.ClassJ: appv1.Class_CLASS_J,
		domain.ClassK: appv1.Class_CLASS_K,
		domain.ClassL: appv1.Class_CLASS_L,
	})
	periodMap = newBimap(map[domain.Period]appv1.Period{
		domain.PeriodPeriod1: appv1.Period_PERIOD_1,
		domain.PeriodPeriod2: appv1.Period_PERIOD_2,
		domain.PeriodPeriod3: appv1.Period_PERIOD_3,
		domain.PeriodPeriod4: appv1.Period_PERIOD_4,
		domain.PeriodPeriod5: appv1.Period_PERIOD_5,
		domain.PeriodPeriod6: appv1.Period_PERIOD_6,
	})
	semesterMap = newBimap(map[domain.CourseSemester]appv1.SubjectSemester{
		domain.CourseSemesterAllYear:         appv1.SubjectSemester_SUBJECT_SEMESTER_ALL_YEAR,
		domain.CourseSemesterH1:              appv1.SubjectSemester_SUBJECT_SEMESTER_H1,
		domain.CourseSemesterH2:              appv1.SubjectSemester_SUBJECT_SEMESTER_H2,
		domain.CourseSemesterQ1:              appv1.SubjectSemester_SUBJECT_SEMESTER_Q1,
		domain.CourseSemesterQ2:              appv1.SubjectSemester_SUBJECT_SEMESTER_Q2,
		domain.CourseSemesterQ3:              appv1.SubjectSemester_SUBJECT_SEMESTER_Q3,
		domain.CourseSemesterQ4:              appv1.SubjectSemester_SUBJECT_SEMESTER_Q4,
		domain.CourseSemesterSummerIntensive: appv1.SubjectSemester_SUBJECT_SEMESTER_SUMMER_INTENSIVE,
		domain.CourseSemesterWinterIntensive: appv1.SubjectSemester_SUBJECT_SEMESTER_WINTER_INTENSIVE,
	})
	classificationMap = newBimap(map[domain.SubjectClassification]appv1.SubjectClassification{
		domain.SubjectClassificationSpecialized:         appv1.SubjectClassification_SUBJECT_CLASSIFICATION_SPECIALIZED,
		domain.SubjectClassificationCultural:            appv1.SubjectClassification_SUBJECT_CLASSIFICATION_CULTURAL,
		domain.SubjectClassificationResearchInstruction: appv1.SubjectClassification_SUBJECT_CLASSIFICATION_RESEARCH_INSTRUCTION,
	})
	culturalCategoryMap = newBimap(map[domain.CulturalSubjectCategory]appv1.SubjectCulturalSubjectCategory{
		domain.CulturalSubjectCategorySociety:       appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_SOCIETY,
		domain.CulturalSubjectCategoryHuman:         appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_HUMAN,
		domain.CulturalSubjectCategoryScience:       appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_SCIENCE,
		domain.CulturalSubjectCategoryHealth:        appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_HEALTH,
		domain.CulturalSubjectCategoryCommunication: appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_COMMUNICATION,
	})
	requirementTypeMap = newBimap(map[domain.SubjectRequirementType]appv1.SubjectRequirementType{
		domain.SubjectRequirementTypeRequired:         appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_REQUIRED,
		domain.SubjectRequirementTypeOptional:         appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL,
		domain.SubjectRequirementTypeOptionalRequired: appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL_REQUIRED,
	})
	floorMap = newBimap(map[domain.Floor]appv1.RoomFloor{
		domain.Floor1:       appv1.RoomFloor_ROOM_FLOOR_1,
		domain.Floor2:       appv1.RoomFloor_ROOM_FLOOR_2,
		domain.Floor3:       appv1.RoomFloor_ROOM_FLOOR_3,
		domain.Floor4:       appv1.RoomFloor_ROOM_FLOOR_4,
		domain.Floor5:       appv1.RoomFloor_ROOM_FLOOR_5,
		domain.Floor6:       appv1.RoomFloor_ROOM_FLOOR_6,
		domain.Floor7:       appv1.RoomFloor_ROOM_FLOOR_7,
		domain.FloorVirtual: appv1.RoomFloor_ROOM_FLOOR_VIRTUAL,
	})
	dayOfWeekMap = newBimap(map[domain.DayOfWeek]dayofweek.DayOfWeek{
		domain.DayOfWeekMonday:    dayofweek.DayOfWeek_MONDAY,
		domain.DayOfWeekTuesday:   dayofweek.DayOfWeek_TUESDAY,
		domain.DayOfWeekWednesday: dayofweek.DayOfWeek_WEDNESDAY,
		domain.DayOfWeekThursday:  dayofweek.DayOfWeek_THURSDAY,
		domain.DayOfWeekFriday:    dayofweek.DayOfWeek_FRIDAY,
		domain.DayOfWeekSaturday:  dayofweek.DayOfWeek_SATURDAY,
		domain.DayOfWeekSunday:    dayofweek.DayOfWeek_SUNDAY,
	})
	calendarStateMap = newBimap(map[domain.PersonalCalendarItemStatus]appv1.PersonalCalendarItem_State{
		domain.PersonalCalendarItemStatusNormal:      appv1.PersonalCalendarItem_STATE_NORMAL,
		domain.PersonalCalendarItemStatusCancelled:   appv1.PersonalCalendarItem_STATE_CANCELLED,
		domain.PersonalCalendarItemStatusMakeup:      appv1.PersonalCalendarItem_STATE_MAKEUP,
		domain.PersonalCalendarItemStatusRoomChanged: appv1.PersonalCalendarItem_STATE_ROOM_CHANGED,
	})
)

// fromProtoEnums は proto の enum のリストを domain の値に変換する。未知の値（UNSPECIFIED を含む）は違反として記録する。
func fromProtoEnums[D comparable, P comparable](m bimap[D, P], values []P, field string, v *violations) []D {
	result := make([]D, 0, len(values))
	for _, p := range values {
		d, ok := m.fromProto[p]
		if !ok {
			v.add(field, "must not contain unspecified or unknown values")
			continue
		}
		result = append(result, d)
	}
	return result
}

func toProtoSubject(s domain.Subject) *appv1.Subject {
	faculties := make([]*appv1.SubjectFaculty, len(s.Faculties))
	for i, f := range s.Faculties {
		faculties[i] = &appv1.SubjectFaculty{Faculty: toProtoFaculty(f.Faculty), IsPrimary: f.IsPrimary}
	}
	return &appv1.Subject{
		Id:        s.ID,
		Name:      s.Name,
		Faculties: faculties,
		Year:      int32(s.Year),
		Semester:  semesterMap.toProto[s.Semester],
		Credit:    int32(s.Credit),
	}
}

// toProtoSubjectDetail は GetSubject 用に詳細のフィールドも設定する。
func toProtoSubjectDetail(s domain.Subject) *appv1.Subject {
	p := toProtoSubject(s)
	p.EligibleAttributes = make([]*appv1.SubjectEligibleAttribute, len(s.EligibleAttributes))
	for i, a := range s.EligibleAttributes {
		attr := &appv1.SubjectEligibleAttribute{Grade: gradeMap.toProto[a.Grade]}
		if a.Class != nil {
			c := classMap.toProto[*a.Class]
			attr.Class = &c
		}
		p.EligibleAttributes[i] = attr
	}
	p.Requirements = make([]*appv1.SubjectRequirement, len(s.Requirements))
	for i, r := range s.Requirements {
		p.Requirements[i] = &appv1.SubjectRequirement{
			Course:          courseMap.toProto[r.Course],
			RequirementType: requirementTypeMap.toProto[r.RequirementType],
		}
	}
	if s.Syllabus != nil {
		p.Syllabus = toProtoSyllabus(*s.Syllabus)
	}
	return p
}

func toProtoSyllabus(s domain.Syllabus) *appv1.SubjectSyllabus {
	return &appv1.SubjectSyllabus{
		Id:                           s.ID,
		Name:                         s.Name,
		EnName:                       s.EnName,
		Grades:                       s.Grades,
		Credit:                       int32(s.Credit),
		FacultyNames:                 s.FacultyNames,
		PracticalHomeFacultyCategory: s.PracticalHomeFacultyCategory,
		MultiplePersonTeachingForm:   s.MultiplePersonTeachingForm,
		TeachingForm:                 s.TeachingForm,
		Summary:                      s.Summary,
		LearningOutcomes:             s.LearningOutcomes,
		Assignments:                  s.Assignments,
		EvaluationMethod:             s.EvaluationMethod,
		Textbooks:                    s.Textbooks,
		ReferenceBooks:               s.ReferenceBooks,
		Prerequisites:                s.Prerequisites,
		PreLearning:                  s.PreLearning,
		PostLearning:                 s.PostLearning,
		Notes:                        s.Notes,
		Keywords:                     s.Keywords,
		TargetCourses:                s.TargetCourses,
		TargetAreas:                  s.TargetAreas,
		Classifications:              s.Classifications,
		TeachingLanguage:             s.TeachingLanguage,
		ContentsAndSchedule:          s.ContentsAndSchedule,
		TeachingAndExamForm:          s.TeachingAndExamForm,
		DsopSubject:                  s.DsopSubject,
	}
}

func toProtoFaculty(f domain.Faculty) *appv1.Faculty {
	return &appv1.Faculty{Id: f.ID, Name: f.Name, Email: f.Email}
}

func toProtoRoom(r domain.Room) *appv1.Room {
	return &appv1.Room{Id: r.ID, Name: r.Name, Floor: floorMap.toProto[r.Floor]}
}

func toProtoRooms(rooms []domain.Room) []*appv1.Room {
	result := make([]*appv1.Room, len(rooms))
	for i, r := range rooms {
		result[i] = toProtoRoom(r)
	}
	return result
}

func mapSlice[T any, U any](items []T, f func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}
