package handler

import (
	"time"
	"uuid"

	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/dayofweek"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const dateLayout = time.DateOnly

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

// parseUUIDs は ID のリストを UUID 文字列にする。不正な ID がある場合は違反として記録する。
func parseUUIDs(ids []string, field string, v *violations) []string {
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		u, err := uuid.Parse(id)
		if err != nil {
			v.add(field, "must contain only valid IDs")
			continue
		}
		result = append(result, u.String())
	}
	return result
}

// enumMap は DB に保存されている文字列と proto の enum の対応。
type enumMap[P comparable] struct {
	toProto   map[string]P
	fromProto map[P]string
}

func newEnumMap[P comparable](pairs map[string]P) enumMap[P] {
	m := enumMap[P]{toProto: pairs, fromProto: make(map[P]string, len(pairs))}
	for s, p := range pairs {
		m.fromProto[p] = s
	}
	return m
}

// fromProtoList は proto の enum のリストを DB の文字列に変換する。未知の値（UNSPECIFIED を含む）は違反として記録する。
func (m enumMap[P]) fromProtoList(values []P, field string, v *violations) []string {
	result := make([]string, 0, len(values))
	for _, p := range values {
		s, ok := m.fromProto[p]
		if !ok {
			v.add(field, "must not contain unspecified or unknown values")
			continue
		}
		result = append(result, s)
	}
	return result
}

func (m enumMap[P]) toProtoPtr(s *string) *P {
	if s == nil {
		return nil
	}
	p, ok := m.toProto[*s]
	if !ok {
		return nil
	}
	return &p
}

var (
	gradeMap = newEnumMap(map[string]appv1.Grade{
		"B1": appv1.Grade_GRADE_B1, "B2": appv1.Grade_GRADE_B2, "B3": appv1.Grade_GRADE_B3, "B4": appv1.Grade_GRADE_B4,
		"M1": appv1.Grade_GRADE_M1, "M2": appv1.Grade_GRADE_M2,
		"D1": appv1.Grade_GRADE_D1, "D2": appv1.Grade_GRADE_D2, "D3": appv1.Grade_GRADE_D3,
	})
	courseMap = newEnumMap(map[string]appv1.Course{
		"InformationSystem": appv1.Course_COURSE_INFORMATION_SYSTEM,
		"InformationDesign": appv1.Course_COURSE_INFORMATION_DESIGN,
		"AdvancedICT":       appv1.Course_COURSE_ADVANCED_ICT,
		"ComplexSystem":     appv1.Course_COURSE_COMPLEX_SYSTEM,
		"IntelligentSystem": appv1.Course_COURSE_INTELLIGENT_SYSTEM,
	})
	classMap = newEnumMap(map[string]appv1.Class{
		"A": appv1.Class_CLASS_A, "B": appv1.Class_CLASS_B, "C": appv1.Class_CLASS_C, "D": appv1.Class_CLASS_D,
		"E": appv1.Class_CLASS_E, "F": appv1.Class_CLASS_F, "G": appv1.Class_CLASS_G, "H": appv1.Class_CLASS_H,
		"I": appv1.Class_CLASS_I, "J": appv1.Class_CLASS_J, "K": appv1.Class_CLASS_K, "L": appv1.Class_CLASS_L,
	})
	periodMap = newEnumMap(map[string]appv1.Period{
		"Period1": appv1.Period_PERIOD_1, "Period2": appv1.Period_PERIOD_2, "Period3": appv1.Period_PERIOD_3,
		"Period4": appv1.Period_PERIOD_4, "Period5": appv1.Period_PERIOD_5, "Period6": appv1.Period_PERIOD_6,
	})
	semesterMap = newEnumMap(map[string]appv1.SubjectSemester{
		"AllYear":         appv1.SubjectSemester_SUBJECT_SEMESTER_ALL_YEAR,
		"H1":              appv1.SubjectSemester_SUBJECT_SEMESTER_H1,
		"H2":              appv1.SubjectSemester_SUBJECT_SEMESTER_H2,
		"Q1":              appv1.SubjectSemester_SUBJECT_SEMESTER_Q1,
		"Q2":              appv1.SubjectSemester_SUBJECT_SEMESTER_Q2,
		"Q3":              appv1.SubjectSemester_SUBJECT_SEMESTER_Q3,
		"Q4":              appv1.SubjectSemester_SUBJECT_SEMESTER_Q4,
		"SummerIntensive": appv1.SubjectSemester_SUBJECT_SEMESTER_SUMMER_INTENSIVE,
		"WinterIntensive": appv1.SubjectSemester_SUBJECT_SEMESTER_WINTER_INTENSIVE,
	})
	classificationMap = newEnumMap(map[string]appv1.SubjectClassification{
		"Specialized":         appv1.SubjectClassification_SUBJECT_CLASSIFICATION_SPECIALIZED,
		"Cultural":            appv1.SubjectClassification_SUBJECT_CLASSIFICATION_CULTURAL,
		"ResearchInstruction": appv1.SubjectClassification_SUBJECT_CLASSIFICATION_RESEARCH_INSTRUCTION,
	})
	culturalCategoryMap = newEnumMap(map[string]appv1.SubjectCulturalSubjectCategory{
		"Society":       appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_SOCIETY,
		"Human":         appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_HUMAN,
		"Science":       appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_SCIENCE,
		"Health":        appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_HEALTH,
		"Communication": appv1.SubjectCulturalSubjectCategory_SUBJECT_CULTURAL_SUBJECT_CATEGORY_COMMUNICATION,
	})
	requirementTypeMap = newEnumMap(map[string]appv1.SubjectRequirementType{
		"Required":         appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_REQUIRED,
		"Optional":         appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL,
		"OptionalRequired": appv1.SubjectRequirementType_SUBJECT_REQUIREMENT_TYPE_OPTIONAL_REQUIRED,
	})
	floorMap = newEnumMap(map[string]appv1.RoomFloor{
		"Floor1": appv1.RoomFloor_ROOM_FLOOR_1, "Floor2": appv1.RoomFloor_ROOM_FLOOR_2,
		"Floor3": appv1.RoomFloor_ROOM_FLOOR_3, "Floor4": appv1.RoomFloor_ROOM_FLOOR_4,
		"Floor5": appv1.RoomFloor_ROOM_FLOOR_5, "Floor6": appv1.RoomFloor_ROOM_FLOOR_6,
		"Floor7": appv1.RoomFloor_ROOM_FLOOR_7, "Virtual": appv1.RoomFloor_ROOM_FLOOR_VIRTUAL,
	})
	dayOfWeekMap = newEnumMap(map[string]dayofweek.DayOfWeek{
		"Monday": dayofweek.DayOfWeek_MONDAY, "Tuesday": dayofweek.DayOfWeek_TUESDAY,
		"Wednesday": dayofweek.DayOfWeek_WEDNESDAY, "Thursday": dayofweek.DayOfWeek_THURSDAY,
		"Friday": dayofweek.DayOfWeek_FRIDAY, "Saturday": dayofweek.DayOfWeek_SATURDAY,
		"Sunday": dayofweek.DayOfWeek_SUNDAY,
	})
)

// subjectPreload は一覧用の Subject に必要な関連を読み込む。prefix は Subject への関連名（例: "Subject."）。
func subjectPreloads(prefix string) []string {
	return []string{prefix + "Faculties.Faculty"}
}

// subjectDetailPreloads は GetSubject 用の関連を読み込む。
func subjectDetailPreloads() []string {
	return []string{"Faculties.Faculty", "EligibleAttributes", "Requirements", "Syllabus"}
}

func toProtoSubject(s *model.Subject) *appv1.Subject {
	if s == nil {
		return nil
	}
	faculties := make([]*appv1.SubjectFaculty, 0, len(s.Faculties))
	for _, f := range s.Faculties {
		faculties = append(faculties, &appv1.SubjectFaculty{Faculty: toProtoFaculty(f.Faculty), IsPrimary: f.IsPrimary})
	}
	return &appv1.Subject{
		Id:        s.ID.String(),
		Name:      s.Name,
		Faculties: faculties,
		Year:      int32(s.Year),
		Semester:  semesterMap.toProto[s.Semester],
		Credit:    int32(s.Credit),
	}
}

// toProtoSubjectDetail は GetSubject 用に詳細のフィールドも設定する。
func toProtoSubjectDetail(s *model.Subject) *appv1.Subject {
	p := toProtoSubject(s)
	p.EligibleAttributes = make([]*appv1.SubjectEligibleAttribute, 0, len(s.EligibleAttributes))
	for _, a := range s.EligibleAttributes {
		p.EligibleAttributes = append(p.EligibleAttributes, &appv1.SubjectEligibleAttribute{
			Grade: gradeMap.toProto[a.Grade],
			Class: classMap.toProtoPtr(a.Class),
		})
	}
	p.Requirements = make([]*appv1.SubjectRequirement, 0, len(s.Requirements))
	for _, r := range s.Requirements {
		p.Requirements = append(p.Requirements, &appv1.SubjectRequirement{
			Course:          courseMap.toProto[r.Course],
			RequirementType: requirementTypeMap.toProto[r.RequirementType],
		})
	}
	if s.Syllabus != nil {
		p.Syllabus = toProtoSyllabus(s.Syllabus)
	}
	return p
}

func toProtoSyllabus(s *model.Syllabus) *appv1.SubjectSyllabus {
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

func toProtoFaculty(f *model.Faculty) *appv1.Faculty {
	if f == nil {
		return nil
	}
	return &appv1.Faculty{Id: f.ID.String(), Name: f.Name, Email: f.Email}
}

func toProtoRoom(r *model.Room) *appv1.Room {
	if r == nil {
		return nil
	}
	return &appv1.Room{Id: r.ID.String(), Name: r.Name, Floor: floorMap.toProto[r.Floor]}
}

func ptrOr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapSlice[T any, U any](items []T, f func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}
