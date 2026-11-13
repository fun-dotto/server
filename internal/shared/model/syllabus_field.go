package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// SyllabusField はシラバスの内容（取得元の値）を表す列。
type SyllabusField struct {
	// Column は DB の列名。
	Column string
	// Label は取得元での表示名。
	Label string
	// Ptr は s の該当フィールドへのポインタを返す。
	Ptr func(s *Syllabus) *string
}

// SyllabusFields はシラバスの内容列の一覧。content_hash の計算と差分の検出に使う。
// 順序を変えると content_hash が変わるため、末尾にのみ追加すること。
var SyllabusFields = []SyllabusField{
	{"display_order", "一覧: 表示順", func(s *Syllabus) *string { return &s.DisplayOrder }},
	{"list_lesson_name", "一覧: 授業科目名", func(s *Syllabus) *string { return &s.ListLessonName }},
	{"list_en_name", "一覧: Course", func(s *Syllabus) *string { return &s.ListEnName }},
	{"list_staff_name", "一覧: 担当教員名", func(s *Syllabus) *string { return &s.ListStaffName }},
	{"list_teaching_form", "一覧: 授業形態", func(s *Syllabus) *string { return &s.ListTeachingForm }},
	{"name", "授業名", func(s *Syllabus) *string { return &s.Name }},
	{"en_name", "Course", func(s *Syllabus) *string { return &s.EnName }},
	{"grades", "配当年次", func(s *Syllabus) *string { return &s.Grades }},
	{"term", "開講期", func(s *Syllabus) *string { return &s.Term }},
	{"credit", "単位数", func(s *Syllabus) *string { return &s.Credit }},
	{"faculty_names", "担当教員名", func(s *Syllabus) *string { return &s.FacultyNames }},
	{"practical_home_faculty_category", "実務家教員区分", func(s *Syllabus) *string { return &s.PracticalHomeFacultyCategory }},
	{"multiple_person_teaching_form", "複数人担当形式", func(s *Syllabus) *string { return &s.MultiplePersonTeachingForm }},
	{"teaching_form", "授業形態", func(s *Syllabus) *string { return &s.TeachingForm }},
	{"summary", "授業の概要", func(s *Syllabus) *string { return &s.Summary }},
	{"learning_outcomes", "授業の到達目標", func(s *Syllabus) *string { return &s.LearningOutcomes }},
	{"assignments", "提出課題等", func(s *Syllabus) *string { return &s.Assignments }},
	{"evaluation_method", "成績の評価方法・基準", func(s *Syllabus) *string { return &s.EvaluationMethod }},
	{"textbooks", "テキスト", func(s *Syllabus) *string { return &s.Textbooks }},
	{"reference_books", "参考書", func(s *Syllabus) *string { return &s.ReferenceBooks }},
	{"prerequisites", "履修条件", func(s *Syllabus) *string { return &s.Prerequisites }},
	{"pre_learning", "事前学習", func(s *Syllabus) *string { return &s.PreLearning }},
	{"post_learning", "事後学習", func(s *Syllabus) *string { return &s.PostLearning }},
	{"notes", "履修上の留意点", func(s *Syllabus) *string { return &s.Notes }},
	{"keywords", "キーワード/keywords", func(s *Syllabus) *string { return &s.Keywords }},
	{"target_courses", "対象コース・領域", func(s *Syllabus) *string { return &s.TargetCourses }},
	{"classifications", "科目群・科目区分", func(s *Syllabus) *string { return &s.Classifications }},
	{"dsop_subject", "DSOP対象科目", func(s *Syllabus) *string { return &s.DsopSubject }},
	{"teaching_language", "教授言語", func(s *Syllabus) *string { return &s.TeachingLanguage }},
	{"teaching_and_exam_form", "授業・試験の形式", func(s *Syllabus) *string { return &s.TeachingAndExamForm }},
	{"contents_and_schedule", "授業内容とスケジュール", func(s *Syllabus) *string { return &s.ContentsAndSchedule }},
	{"en_summary", "(E)Course outline", func(s *Syllabus) *string { return &s.EnSummary }},
	{"en_learning_outcomes", "(E)Course Objectives", func(s *Syllabus) *string { return &s.EnLearningOutcomes }},
	{"en_contents_and_schedule", "(E)Course Schedule", func(s *Syllabus) *string { return &s.EnContentsAndSchedule }},
	{"en_evaluation_method", "(E)Assessment", func(s *Syllabus) *string { return &s.EnEvaluationMethod }},
	{"en_teaching_language", "(E)Language of Instruction", func(s *Syllabus) *string { return &s.EnTeachingLanguage }},
	{"en_prerequisites", "(E)Requirements for registration", func(s *Syllabus) *string { return &s.EnPrerequisites }},
	{"en_pre_post_learning", "(E)Prior/Post Assignment", func(s *Syllabus) *string { return &s.EnPrePostLearning }},
	{"en_textbooks", "(E)Textbooks", func(s *Syllabus) *string { return &s.EnTextbooks }},
	{"en_notes", "(E)Note", func(s *Syllabus) *string { return &s.EnNotes }},
	{"en_teaching_and_exam_form", "(E)Type of class and exam", func(s *Syllabus) *string { return &s.EnTeachingAndExamForm }},
	{"diploma_policies", "ディプロマポリシー", func(s *Syllabus) *string { return &s.DiplomaPolicies }},
	{"diploma_policy_correspondence", "ディプロマポリシーとの対応", func(s *Syllabus) *string { return &s.DiplomaPolicyCorrespondence }},
}

// ComputeContentHash は内容列から content_hash を計算する。
func (s *Syllabus) ComputeContentHash() string {
	values := make([]string, len(SyllabusFields))
	for i, f := range SyllabusFields {
		values[i] = *f.Ptr(s)
	}
	b, err := json.Marshal(values)
	if err != nil {
		// []string の Marshal は失敗しない
		panic(err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// ChangedSyllabusFields は a と b で値が異なる内容列を返す。
func ChangedSyllabusFields(a, b *Syllabus) []SyllabusField {
	var changed []SyllabusField
	for _, f := range SyllabusFields {
		if *f.Ptr(a) != *f.Ptr(b) {
			changed = append(changed, f)
		}
	}
	return changed
}
