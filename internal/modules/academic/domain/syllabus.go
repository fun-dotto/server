package domain

import "time"

// Syllabus は教務システムのシラバス（取得元の写し）。
type Syllabus struct {
	ID                           string
	LessonID                     string
	Year                         int
	DisplayOrder                 string
	ListLessonName               string
	ListEnName                   string
	ListStaffName                string
	ListTeachingForm             string
	Name                         string
	EnName                       string
	Grades                       string
	Term                         string
	Credit                       string
	FacultyNames                 string
	PracticalHomeFacultyCategory string
	MultiplePersonTeachingForm   string
	TeachingForm                 string
	Summary                      string
	LearningOutcomes             string
	Assignments                  string
	EvaluationMethod             string
	Textbooks                    string
	ReferenceBooks               string
	Prerequisites                string
	PreLearning                  string
	PostLearning                 string
	Notes                        string
	Keywords                     string
	TargetCourses                string
	Classifications              string
	DsopSubject                  string
	TeachingLanguage             string
	TeachingAndExamForm          string
	ContentsAndSchedule          string
	EnSummary                    string
	EnLearningOutcomes           string
	EnContentsAndSchedule        string
	EnEvaluationMethod           string
	EnTeachingLanguage           string
	EnPrerequisites              string
	EnPrePostLearning            string
	EnTextbooks                  string
	EnNotes                      string
	EnTeachingAndExamForm        string
	DiplomaPolicies              string
	DiplomaPolicyCorrespondence  string
	ContentHash                  string
	FetchedAt                    time.Time
}
