package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUsersV1Detail_ProxiesUserAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user":{"id":"user-1","name":"Jane Doe","email":"jane@example.com"}}`))
	}))
	defer server.Close()

	h := newTestHandler(t, server.URL)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/users/user-1", nil)
	setAdminClaim(c)

	h.UsersV1Detail(c, "user-1")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if gotPath != "/v1/users/user-1" {
		t.Fatalf("path = %q, want %q", gotPath, "/v1/users/user-1")
	}

	var body struct {
		User struct {
			Id string `json:"id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.User.Id != "user-1" {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
}
