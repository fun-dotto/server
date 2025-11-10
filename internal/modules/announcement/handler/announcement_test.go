package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewHandler()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.AnnouncementsList(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var announcements []api.Announcement
	json.Unmarshal(w.Body.Bytes(), &announcements)

	assert.NotEmpty(t, announcements)
	assert.Equal(t, "1", announcements[0].Id)
}
