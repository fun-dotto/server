package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	academic_api "github.com/fun-dotto/server/gen/academic"
	"github.com/gin-gonic/gin"
)

func TestCourseRegistrationsV1Create_ProxiesAcademicAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotMethod string
	var gotPath string
	var gotBody academic_api.CourseRegistrationRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path

		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"courseRegistration":{"id":"reg-1","userId":"user-1","subject":{"id":"subject-1","name":"Algorithms","faculties":[]}}}`))
	}))
	defer server.Close()

	h := newTestHandler(t, server.URL)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(
		http.MethodPost,
		"/v1/courseRegistrations",
		bytes.NewBufferString(`{"subjectId":"subject-1","userId":"user-1"}`),
	)
	c.Request.Header.Set("Content-Type", "application/json")
	setAdminClaim(c)

	h.CourseRegistrationsV1Create(c)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if gotMethod != http.MethodPost || gotPath != "/v1/courseRegistrations" {
		t.Fatalf("method/path = %s %s", gotMethod, gotPath)
	}
	if gotBody.SubjectId != "subject-1" || gotBody.UserId != "user-1" {
		t.Fatalf("unexpected upstream request body: %+v", gotBody)
	}
}
