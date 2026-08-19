package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/fun-dotto/server/gen/admin"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func TestPersonalCalendarItemsV1List_ProxiesAcademicAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	wantDate := openapi_types.Date{Time: time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC)}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("userId"); got != "user-1" {
			t.Fatalf("userId = %q, want %q", got, "user-1")
		}
		if got := query.Get("dates"); got == "" {
			t.Fatal("dates query parameter is empty")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"personalCalendarItems":[{"date":"2026-03-26","slot":{"dayOfWeek":"Thursday","period":"first"},"timetableItem":{"id":"item-1","subject":{"id":"subject-1","name":"Algorithms"},"slot":{"dayOfWeek":"Thursday","period":"first"},"rooms":[]}}]}`))
	}))
	defer server.Close()

	h := newTestHandler(t, server.URL)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/personalCalendarItems", nil)
	setAdminClaim(c)

	h.PersonalCalendarItemsV1List(c, api.PersonalCalendarItemsV1ListParams{
		UserId: "user-1",
		Dates:  []openapi_types.Date{wantDate},
	})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		PersonalCalendarItems []struct {
			Date string `json:"date"`
		} `json:"personalCalendarItems"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.PersonalCalendarItems) != 1 || body.PersonalCalendarItems[0].Date != "2026-03-26" {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
}
