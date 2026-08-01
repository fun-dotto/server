package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	api "github.com/fun-dotto/server/gen/admin"
	"github.com/gin-gonic/gin"
)

func TestSubjectsV1List_UsesRenamedQueryParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"subjects":[]}`))
	}))
	defer server.Close()

	h := newTestHandler(t, server.URL)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/subjects", nil)
	setAdminClaim(c)

	q := "math"
	year := 2026
	grades := []api.DottoFoundationV1Grade{api.B1}
	courses := []api.DottoFoundationV1Course{api.InformationSystem}
	classes := []api.DottoFoundationV1Class{api.A}
	classifications := []api.DottoFoundationV1SubjectClassification{api.Cultural}
	semesters := []api.DottoFoundationV1CourseSemester{api.Q1, api.Q2}
	requirementTypes := []api.DottoFoundationV1SubjectRequirementType{api.Optional}
	categories := []api.DottoFoundationV1CulturalSubjectCategory{api.Society}

	h.SubjectsV1List(c, api.SubjectsV1ListParams{
		Q:                         &q,
		Grades:                    &grades,
		Courses:                   &courses,
		Classes:                   &classes,
		Classifications:           &classifications,
		Year:                      &year,
		Semesters:                 &semesters,
		RequirementTypes:          &requirementTypes,
		CulturalSubjectCategories: &categories,
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if gotQuery.Get("grades") == "" || gotQuery.Get("grade") != "" {
		t.Fatalf("query grades = %q, grade = %q", gotQuery.Get("grades"), gotQuery.Get("grade"))
	}
	if gotQuery.Get("classes") == "" || gotQuery.Get("class") != "" {
		t.Fatalf("query classes = %q, class = %q", gotQuery.Get("classes"), gotQuery.Get("class"))
	}
	if gotQuery.Get("classifications") == "" || gotQuery.Get("classification") != "" {
		t.Fatalf("query classifications = %q, classification = %q", gotQuery.Get("classifications"), gotQuery.Get("classification"))
	}
	if gotQuery.Get("semesters") == "" || gotQuery.Get("semester") != "" {
		t.Fatalf("query semesters = %q, semester = %q", gotQuery.Get("semesters"), gotQuery.Get("semester"))
	}
	if gotQuery.Get("requirementTypes") == "" || gotQuery.Get("requirementType") != "" {
		t.Fatalf("query requirementTypes = %q, requirementType = %q", gotQuery.Get("requirementTypes"), gotQuery.Get("requirementType"))
	}
	if gotQuery.Get("culturalSubjectCategories") == "" || gotQuery.Get("culturalSubjectCategory") != "" {
		t.Fatalf("query culturalSubjectCategories = %q, culturalSubjectCategory = %q", gotQuery.Get("culturalSubjectCategories"), gotQuery.Get("culturalSubjectCategory"))
	}
}
