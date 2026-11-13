package handler

import (
	"testing"

	"github.com/fun-dotto/server/internal/modules/academic/domain"
)

func TestSyllabusToAPI(t *testing.T) {
	input := domain.Syllabus{
		ID:                           "0f8fad5b-d9cb-469f-a165-70867728950e",
		LessonID:                     "100101",
		Year:                         2026,
		Term:                         "前期",
		EnSummary:                    "Basics of algorithms",
		DiplomaPolicyCorrespondence:  "DP2,DP4",
		Name:                         "アルゴリズム論",
		EnName:                       "Algorithm Theory",
		Grades:                       "B2",
		Credit:                       "2単位",
		FacultyNames:                 "田中太郎",
		PracticalHomeFacultyCategory: "該当しない",
		MultiplePersonTeachingForm:   "オムニバス",
		TeachingForm:                 "講義",
		Summary:                      "アルゴリズムの基礎を学ぶ",
		LearningOutcomes:             "基本的なアルゴリズムを理解する",
		Assignments:                  "毎週のレポート",
		EvaluationMethod:             "期末試験60%、レポート40%",
		Textbooks:                    "アルゴリズムとデータ構造",
		ReferenceBooks:               "Introduction to Algorithms",
		Prerequisites:                "プログラミング基礎",
		PreLearning:                  "教科書の該当箇所を読む",
		PostLearning:                 "演習問題を解く",
		Notes:                        "ノートPC持参",
		Keywords:                     "アルゴリズム,データ構造,計算量",
		TargetCourses:                "情報システムコース",
		Classifications:              "専門科目",
		TeachingLanguage:             "日本語",
		ContentsAndSchedule:          "第1回:ガイダンス",
		TeachingAndExamForm:          "対面",
		DsopSubject:                  "対象外",
	}

	got := syllabusToAPI(input)

	checks := []struct {
		field string
		got   string
		want  string
	}{
		{"Id", got.Id, input.ID},
		{"Name", got.Name, input.Name},
		{"EnName", got.EnName, input.EnName},
		{"Grades", got.Grades, input.Grades},
		{"FacultyNames", got.FacultyNames, input.FacultyNames},
		{"PracticalHomeFacultyCategory", got.PracticalHomeFacultyCategory, input.PracticalHomeFacultyCategory},
		{"MultiplePersonTeachingForm", got.MultiplePersonTeachingForm, input.MultiplePersonTeachingForm},
		{"TeachingForm", got.TeachingForm, input.TeachingForm},
		{"Summary", got.Summary, input.Summary},
		{"LearningOutcomes", got.LearningOutcomes, input.LearningOutcomes},
		{"Assignments", got.Assignments, input.Assignments},
		{"EvaluationMethod", got.EvaluationMethod, input.EvaluationMethod},
		{"Textbooks", got.Textbooks, input.Textbooks},
		{"ReferenceBooks", got.ReferenceBooks, input.ReferenceBooks},
		{"Prerequisites", got.Prerequisites, input.Prerequisites},
		{"PreLearning", got.PreLearning, input.PreLearning},
		{"PostLearning", got.PostLearning, input.PostLearning},
		{"Notes", got.Notes, input.Notes},
		{"Keywords", got.Keywords, input.Keywords},
		{"TargetCourses", got.TargetCourses, input.TargetCourses},
		{"LessonId", got.LessonId, input.LessonID},
		{"Term", got.Term, input.Term},
		{"Credit", got.Credit, input.Credit},
		{"EnSummary", got.EnSummary, input.EnSummary},
		{"DiplomaPolicyCorrespondence", got.DiplomaPolicyCorrespondence, input.DiplomaPolicyCorrespondence},
		{"Classifications", got.Classifications, input.Classifications},
		{"TeachingLanguage", got.TeachingLanguage, input.TeachingLanguage},
		{"ContentsAndSchedule", got.ContentsAndSchedule, input.ContentsAndSchedule},
		{"TeachingAndExamForm", got.TeachingAndExamForm, input.TeachingAndExamForm},
		{"DsopSubject", got.DsopSubject, input.DsopSubject},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s: got %q, want %q", c.field, c.got, c.want)
		}
	}

	if got.Year != input.Year {
		t.Errorf("Year: got %d, want %d", got.Year, input.Year)
	}
}
