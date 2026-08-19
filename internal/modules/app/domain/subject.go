package domain

// Subject 科目
type Subject struct {
	ID                 string
	Name               string
	Year               int
	Credit             int
	Semester           CourseSemester
	Faculties          []SubjectFaculty
	Requirements       []SubjectRequirement
	EligibleAttributes []SubjectTargetClass
	Syllabus           *Syllabus
}

// SubjectFaculty 科目担当教員
type SubjectFaculty struct {
	Faculty   Faculty
	IsPrimary bool
}

// Faculty 教員
type Faculty struct {
	ID    string
	Name  string
	Email string
}

// SubjectRequirement 科目の必修/選択区分
type SubjectRequirement struct {
	Course          Course
	RequirementType SubjectRequirementType
}

// SubjectTargetClass 対象学年・クラス
type SubjectTargetClass struct {
	Grade Grade
	Class *Class
}

// Syllabus シラバス
type Syllabus struct {
	ID                           string
	Name                         string
	EnName                       string
	Summary                      string
	LearningOutcomes             string
	ContentsAndSchedule          string
	PreLearning                  string
	PostLearning                 string
	Assignments                  string
	EvaluationMethod             string
	Textbooks                    string
	ReferenceBooks               string
	Prerequisites                string
	Notes                        string
	Keywords                     string
	Classifications              string
	Grades                       string
	Credit                       int
	FacultyNames                 string
	TeachingForm                 string
	TeachingAndExamForm          string
	TeachingLanguage             string
	MultiplePersonTeachingForm   string
	PracticalHomeFacultyCategory string
	DsopSubject                  string
	TargetAreas                  string
	TargetCourses                string
}

// CourseSemester 開講時期
type CourseSemester string

const (
	CourseSemesterQ1              CourseSemester = "Q1"
	CourseSemesterQ2              CourseSemester = "Q2"
	CourseSemesterQ3              CourseSemester = "Q3"
	CourseSemesterQ4              CourseSemester = "Q4"
	CourseSemesterH1              CourseSemester = "H1"
	CourseSemesterH2              CourseSemester = "H2"
	CourseSemesterAllYear         CourseSemester = "AllYear"
	CourseSemesterSummerIntensive CourseSemester = "SummerIntensive"
	CourseSemesterWinterIntensive CourseSemester = "WinterIntensive"
)

// Grade 学年
type Grade string

const (
	GradeB1 Grade = "B1"
	GradeB2 Grade = "B2"
	GradeB3 Grade = "B3"
	GradeB4 Grade = "B4"
	GradeM1 Grade = "M1"
	GradeM2 Grade = "M2"
	GradeD1 Grade = "D1"
	GradeD2 Grade = "D2"
	GradeD3 Grade = "D3"
)

// Class クラス
type Class string

const (
	ClassA Class = "A"
	ClassB Class = "B"
	ClassC Class = "C"
	ClassD Class = "D"
	ClassE Class = "E"
	ClassF Class = "F"
	ClassG Class = "G"
	ClassH Class = "H"
	ClassI Class = "I"
	ClassJ Class = "J"
	ClassK Class = "K"
	ClassL Class = "L"
)

// Course コース
type Course string

const (
	CourseInformationSystem Course = "InformationSystem" // 情報システムコース
	CourseInformationDesign Course = "InformationDesign" // 情報デザインコース
	CourseComplexSystem     Course = "ComplexSystem"     // 複雑系コース
	CourseIntelligentSystem Course = "IntelligentSystem" // 知能システムコース
	CourseAdvancedICT       Course = "AdvancedICT"       // 高度ICTコース
)

// SubjectRequirementType 必修・選択区分
type SubjectRequirementType string

const (
	SubjectRequirementTypeRequired         SubjectRequirementType = "Required"         // 必修
	SubjectRequirementTypeOptional         SubjectRequirementType = "Optional"         // 選択
	SubjectRequirementTypeOptionalRequired SubjectRequirementType = "OptionalRequired" // 選択必修
)

// SubjectClassification 科目分類
type SubjectClassification string

const (
	SubjectClassificationSpecialized         SubjectClassification = "Specialized"         // 専門科目
	SubjectClassificationCultural            SubjectClassification = "Cultural"            // 教養科目
	SubjectClassificationResearchInstruction SubjectClassification = "ResearchInstruction" // 研究指導科目
)

// CulturalSubjectCategory 教養科目カテゴリ
type CulturalSubjectCategory string

const (
	CulturalSubjectCategoryHuman         CulturalSubjectCategory = "Human"         // 人間の形成
	CulturalSubjectCategorySociety       CulturalSubjectCategory = "Society"       // 社会への参加
	CulturalSubjectCategoryScience       CulturalSubjectCategory = "Science"       // 科学技術と環境の理解
	CulturalSubjectCategoryHealth        CulturalSubjectCategory = "Health"        // 健康の保持
	CulturalSubjectCategoryCommunication CulturalSubjectCategory = "Communication" // コミュニケーション
)

// SubjectQuery 科目検索クエリ
type SubjectQuery struct {
	Q                       string
	Grade                   []Grade
	Courses                 []Course
	Class                   []Class
	Classification          []SubjectClassification
	Semester                []CourseSemester
	RequirementType         []SubjectRequirementType
	CulturalSubjectCategory []CulturalSubjectCategory
}
