package domain

type CourseSemester string

const (
	CourseSemesterAllYear         CourseSemester = "AllYear"
	CourseSemesterH1              CourseSemester = "H1"
	CourseSemesterH2              CourseSemester = "H2"
	CourseSemesterQ1              CourseSemester = "Q1"
	CourseSemesterQ2              CourseSemester = "Q2"
	CourseSemesterQ3              CourseSemester = "Q3"
	CourseSemesterQ4              CourseSemester = "Q4"
	CourseSemesterSummerIntensive CourseSemester = "SummerIntensive"
	CourseSemesterWinterIntensive CourseSemester = "WinterIntensive"
)

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

type SubjectRequirementType string

const (
	SubjectRequirementTypeRequired         SubjectRequirementType = "Required"
	SubjectRequirementTypeOptional         SubjectRequirementType = "Optional"
	SubjectRequirementTypeOptionalRequired SubjectRequirementType = "OptionalRequired"
)

type CourseType string

const (
	CourseTypeInformationSystem CourseType = "InformationSystem"
	CourseTypeInformationDesign CourseType = "InformationDesign"
	CourseTypeAdvancedICT       CourseType = "AdvancedICT"
	CourseTypeComplexSystem     CourseType = "ComplexSystem"
	CourseTypeIntelligentSystem CourseType = "IntelligentSystem"
)

type SubjectClassification string

const (
	SubjectClassificationSpecialized         SubjectClassification = "Specialized"
	SubjectClassificationCultural            SubjectClassification = "Cultural"
	SubjectClassificationResearchInstruction SubjectClassification = "ResearchInstruction"
)

type CulturalSubjectCategory string

const (
	CulturalSubjectCategorySociety       CulturalSubjectCategory = "Society"
	CulturalSubjectCategoryHuman         CulturalSubjectCategory = "Human"
	CulturalSubjectCategoryScience       CulturalSubjectCategory = "Science"
	CulturalSubjectCategoryHealth        CulturalSubjectCategory = "Health"
	CulturalSubjectCategoryCommunication CulturalSubjectCategory = "Communication"
)
