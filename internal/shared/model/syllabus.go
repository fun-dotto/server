package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// Syllabus は教務システムのシラバスの取得元の写し。値は加工せず文字列のまま保持する。
// 取得した内容が前回と異なる場合は新しい行（バージョン）として追加するため、
// 同じ (LessonID, Year) に複数のバージョンが存在しうる。
type Syllabus struct {
	Common

	LessonID string `gorm:"not null;uniqueIndex:idx_syllabuses_version,priority:1;index:idx_syllabuses_lesson_year_fetched_at,priority:1"`
	Year     int    `gorm:"not null;uniqueIndex:idx_syllabuses_version,priority:2;index:idx_syllabuses_lesson_year_fetched_at,priority:2"`

	DisplayOrder     string `gorm:"not null"`
	ListLessonName   string `gorm:"not null"`
	ListEnName       string `gorm:"not null"`
	ListStaffName    string `gorm:"not null"`
	ListTeachingForm string `gorm:"not null"`

	Name                         string `gorm:"not null"`
	EnName                       string `gorm:"not null"`
	Grades                       string `gorm:"not null"`
	Term                         string `gorm:"not null"`
	Credit                       string `gorm:"not null"`
	FacultyNames                 string `gorm:"not null"`
	PracticalHomeFacultyCategory string `gorm:"not null"`
	MultiplePersonTeachingForm   string `gorm:"not null"`
	TeachingForm                 string `gorm:"not null"`
	Summary                      string `gorm:"not null"`
	LearningOutcomes             string `gorm:"not null"`
	Assignments                  string `gorm:"not null"`
	EvaluationMethod             string `gorm:"not null"`
	Textbooks                    string `gorm:"not null"`
	ReferenceBooks               string `gorm:"not null"`
	Prerequisites                string `gorm:"not null"`
	PreLearning                  string `gorm:"not null"`
	PostLearning                 string `gorm:"not null"`
	Notes                        string `gorm:"not null"`
	Keywords                     string `gorm:"not null"`
	TargetCourses                string `gorm:"not null"`
	Classifications              string `gorm:"not null"`
	DsopSubject                  string `gorm:"not null"`
	TeachingLanguage             string `gorm:"not null"`
	TeachingAndExamForm          string `gorm:"not null"`
	ContentsAndSchedule          string `gorm:"not null"`

	EnSummary             string `gorm:"not null"`
	EnLearningOutcomes    string `gorm:"not null"`
	EnContentsAndSchedule string `gorm:"not null"`
	EnEvaluationMethod    string `gorm:"not null"`
	EnTeachingLanguage    string `gorm:"not null"`
	EnPrerequisites       string `gorm:"not null"`
	EnPrePostLearning     string `gorm:"not null"`
	EnTextbooks           string `gorm:"not null"`
	EnNotes               string `gorm:"not null"`
	EnTeachingAndExamForm string `gorm:"not null"`

	DiplomaPolicies             string `gorm:"not null"`
	DiplomaPolicyCorrespondence string `gorm:"not null"`

	// 取得値から計算したハッシュ。
	ContentHash string `gorm:"not null;uniqueIndex:idx_syllabuses_version,priority:3"`
	// 取得日時。同一内容を再取得した場合は最終確認日時として更新する。
	FetchedAt time.Time `gorm:"not null;default:current_timestamp;index:idx_syllabuses_lesson_year_fetched_at,priority:3"`
}

// ComputeContentHash は取得元の値（lesson_id・year と各項目）から内容のハッシュを計算する。
// 内容が同じであれば同じ値になるため、新しいバージョンかどうかの判定に使う。
func (s *Syllabus) ComputeContentHash() string {
	h := sha256.New()
	write := func(v string) {
		// 区切りの曖昧さをなくすため長さを前置する
		fmt.Fprintf(h, "%d:%s;", len(v), v)
	}
	write(s.LessonID)
	write(strconv.Itoa(s.Year))
	write(s.DisplayOrder)
	write(s.ListLessonName)
	write(s.ListEnName)
	write(s.ListStaffName)
	write(s.ListTeachingForm)
	write(s.Name)
	write(s.EnName)
	write(s.Grades)
	write(s.Term)
	write(s.Credit)
	write(s.FacultyNames)
	write(s.PracticalHomeFacultyCategory)
	write(s.MultiplePersonTeachingForm)
	write(s.TeachingForm)
	write(s.Summary)
	write(s.LearningOutcomes)
	write(s.Assignments)
	write(s.EvaluationMethod)
	write(s.Textbooks)
	write(s.ReferenceBooks)
	write(s.Prerequisites)
	write(s.PreLearning)
	write(s.PostLearning)
	write(s.Notes)
	write(s.Keywords)
	write(s.TargetCourses)
	write(s.Classifications)
	write(s.DsopSubject)
	write(s.TeachingLanguage)
	write(s.TeachingAndExamForm)
	write(s.ContentsAndSchedule)
	write(s.EnSummary)
	write(s.EnLearningOutcomes)
	write(s.EnContentsAndSchedule)
	write(s.EnEvaluationMethod)
	write(s.EnTeachingLanguage)
	write(s.EnPrerequisites)
	write(s.EnPrePostLearning)
	write(s.EnTextbooks)
	write(s.EnNotes)
	write(s.EnTeachingAndExamForm)
	write(s.DiplomaPolicies)
	write(s.DiplomaPolicyCorrespondence)
	return hex.EncodeToString(h.Sum(nil))
}

// SubjectOrderByLessonID は科目を紐付いたシラバスの授業コード順に並べる ORDER BY 句。
// subjects テーブルを参照するクエリで使う。
const SubjectOrderByLessonID = "(SELECT syllabuses.lesson_id FROM syllabuses WHERE syllabuses.id = subjects.syllabus_id) ASC"
