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

	tests := []struct {
		name         string
		setupContext func(c *gin.Context)
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:         "正常にお知らせ一覧が取得できる",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err, "JSONのパースに失敗しました")
				assert.NotEmpty(t, announcements, "アナウンスメントが空です")
			},
		},
		{
			name:         "Content-Typeがapplication/jsonである",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			},
		},
		{
			name:         "レスポンスが配列形式である",
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result interface{}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				_, isArray := result.([]interface{})
				assert.True(t, isArray, "レスポンスが配列形式ではありません")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			h.AnnouncementsList(c)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
