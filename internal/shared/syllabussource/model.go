package syllabussource

import (
	"slices"

	"github.com/fun-dotto/server/internal/shared/model"
)

// detailFields は詳細ページの見出し（NormalizeLabel 済み）と保存先の列の対応。
var detailFields = map[string]func(*model.Syllabus) *string{
	"授業名":                              func(s *model.Syllabus) *string { return &s.Name },
	"Course":                           func(s *model.Syllabus) *string { return &s.EnName },
	"配当年次":                             func(s *model.Syllabus) *string { return &s.Grades },
	"開講期":                              func(s *model.Syllabus) *string { return &s.Term },
	"予定開講期":                            func(s *model.Syllabus) *string { return &s.Term },
	"単位数":                              func(s *model.Syllabus) *string { return &s.Credit },
	"担当教員名":                            func(s *model.Syllabus) *string { return &s.FacultyNames },
	"実務家教員区分":                          func(s *model.Syllabus) *string { return &s.PracticalHomeFacultyCategory },
	"複数人担当形式":                          func(s *model.Syllabus) *string { return &s.MultiplePersonTeachingForm },
	"授業形態":                             func(s *model.Syllabus) *string { return &s.TeachingForm },
	"授業の概要":                            func(s *model.Syllabus) *string { return &s.Summary },
	"授業の到達目標":                          func(s *model.Syllabus) *string { return &s.LearningOutcomes },
	"提出課題等":                            func(s *model.Syllabus) *string { return &s.Assignments },
	"成績の評価方法・基準":                       func(s *model.Syllabus) *string { return &s.EvaluationMethod },
	"テキスト":                             func(s *model.Syllabus) *string { return &s.Textbooks },
	"参考書":                              func(s *model.Syllabus) *string { return &s.ReferenceBooks },
	"履修条件":                             func(s *model.Syllabus) *string { return &s.Prerequisites },
	"事前学習":                             func(s *model.Syllabus) *string { return &s.PreLearning },
	"事後学習":                             func(s *model.Syllabus) *string { return &s.PostLearning },
	"履修上の留意点":                          func(s *model.Syllabus) *string { return &s.Notes },
	"キーワード/keywords":                   func(s *model.Syllabus) *string { return &s.Keywords },
	"対象コース・領域":                         func(s *model.Syllabus) *string { return &s.TargetCourses },
	"科目群・科目区分":                         func(s *model.Syllabus) *string { return &s.Classifications },
	"DSOP対象科目":                         func(s *model.Syllabus) *string { return &s.DsopSubject },
	"教授言語":                             func(s *model.Syllabus) *string { return &s.TeachingLanguage },
	"授業・試験の形式":                         func(s *model.Syllabus) *string { return &s.TeachingAndExamForm },
	"授業内容とスケジュール":                      func(s *model.Syllabus) *string { return &s.ContentsAndSchedule },
	"(E)Course outline":                func(s *model.Syllabus) *string { return &s.EnSummary },
	"(E)Course Objectives":             func(s *model.Syllabus) *string { return &s.EnLearningOutcomes },
	"(E)Course Schedule":               func(s *model.Syllabus) *string { return &s.EnContentsAndSchedule },
	"(E)Assessment":                    func(s *model.Syllabus) *string { return &s.EnEvaluationMethod },
	"(E)Language of Instruction":       func(s *model.Syllabus) *string { return &s.EnTeachingLanguage },
	"(E)Requirements for registration": func(s *model.Syllabus) *string { return &s.EnPrerequisites },
	"(E)Prior/Post Assignment":         func(s *model.Syllabus) *string { return &s.EnPrePostLearning },
	"(E)Textbooks":                     func(s *model.Syllabus) *string { return &s.EnTextbooks },
	"(E)Note":                          func(s *model.Syllabus) *string { return &s.EnNotes },
	"(E)Type of class and exam":        func(s *model.Syllabus) *string { return &s.EnTeachingAndExamForm },
	"ディプロマポリシー":                        func(s *model.Syllabus) *string { return &s.DiplomaPolicies },
	"ディプロマポリシーとの対応":                    func(s *model.Syllabus) *string { return &s.DiplomaPolicyCorrespondence },
}

// ToModel は一覧の行と詳細ページの値から syllabuses の行を作る。
// 対応する列がない見出しを unknown として返す（取得元に項目が追加されたことの検知用）。
func ToModel(year int, e ListEntry, d Detail) (m model.Syllabus, unknown []string) {
	m = model.Syllabus{
		LessonID:         e.LessonID,
		Year:             year,
		DisplayOrder:     e.DisplayOrder,
		ListLessonName:   e.LessonName,
		ListEnName:       e.EnName,
		ListStaffName:    e.StaffName,
		ListTeachingForm: e.TeachingForm,
	}
	for label, value := range d {
		ptr, ok := detailFields[label]
		if !ok {
			unknown = append(unknown, label)
			continue
		}
		*ptr(&m) = value
	}
	slices.Sort(unknown)
	return m, unknown
}
