package model

import (
	"time"
	"uuid"
)

// Syllabus は教務システムのシラバス（一覧・詳細ページ）の写し。
// 値は加工せず文字列のまま保存する（空欄・非表示は空文字）。
// 取得した内容が前回と異なる場合は新しい行（バージョン）として追加する。
type Syllabus struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LessonID string    `gorm:"not null;uniqueIndex:idx_syllabuses_lesson_id_year_content_hash,priority:1;index:idx_syllabuses_lesson_id_year_fetched_at,priority:1"`
	Year     int       `gorm:"not null;uniqueIndex:idx_syllabuses_lesson_id_year_content_hash,priority:2;index:idx_syllabuses_lesson_id_year_fetched_at,priority:2"`

	// 一覧ページの値
	DisplayOrder     string `gorm:"not null;default:''"`
	ListLessonName   string `gorm:"not null;default:''"`
	ListEnName       string `gorm:"not null;default:''"`
	ListStaffName    string `gorm:"not null;default:''"`
	ListTeachingForm string `gorm:"not null;default:''"`

	// 詳細ページの値
	Name                         string `gorm:"not null"`
	EnName                       string `gorm:"not null"`
	Grades                       string `gorm:"not null"`
	Term                         string `gorm:"not null;default:''"`
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
	EnSummary                    string `gorm:"not null;default:''"`
	EnLearningOutcomes           string `gorm:"not null;default:''"`
	EnContentsAndSchedule        string `gorm:"not null;default:''"`
	EnEvaluationMethod           string `gorm:"not null;default:''"`
	EnTeachingLanguage           string `gorm:"not null;default:''"`
	EnPrerequisites              string `gorm:"not null;default:''"`
	EnPrePostLearning            string `gorm:"not null;default:''"`
	EnTextbooks                  string `gorm:"not null;default:''"`
	EnNotes                      string `gorm:"not null;default:''"`
	EnTeachingAndExamForm        string `gorm:"not null;default:''"`
	DiplomaPolicies              string `gorm:"not null;default:''"`
	DiplomaPolicyCorrespondence  string `gorm:"not null;default:''"`

	// ContentHash は取得値から計算したハッシュ。同じ (lesson_id, year) で内容が同じかを判定する。
	ContentHash string `gorm:"not null;uniqueIndex:idx_syllabuses_lesson_id_year_content_hash,priority:3"`
	// FetchedAt は取得日時。同一内容を再取得した場合は最終確認日時として更新する。
	FetchedAt time.Time `gorm:"not null;default:current_timestamp;index:idx_syllabuses_lesson_id_year_fetched_at,priority:3"`

	CreatedAt time.Time `gorm:"autoCreateTime;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;default:current_timestamp"`
}
