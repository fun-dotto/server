package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/repository"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		isActive     *bool
		setupContext func(c *gin.Context)
		wantCode     int
		validate     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:         "正常にお知らせ一覧が取得できる",
			isActive:     boolPtr(true),
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
			name:         "isActive=trueで有効なお知らせのみ取得できる",
			isActive:     boolPtr(true),
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err)
				assert.Len(t, announcements, 1, "有効なお知らせは1件のはずです")
				assert.True(t, announcements[0].IsActive, "IsActiveがtrueではありません")
			},
		},
		{
			name:         "isActive=falseで無効なお知らせのみ取得できる",
			isActive:     boolPtr(false),
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err)
				assert.Len(t, announcements, 1, "無効なお知らせは1件のはずです")
				assert.False(t, announcements[0].IsActive, "IsActiveがfalseではありません")
			},
		},
		{
			name:         "isActive=nilで全件取得できる",
			isActive:     nil,
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err)
				assert.Len(t, announcements, 2, "全件（2件）取得できるはずです")
			},
		},
		{
			name:         "Content-Typeがapplication/jsonである",
			isActive:     boolPtr(true),
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
			},
		},
		{
			name:         "レスポンスが配列形式である",
			isActive:     boolPtr(true),
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
			mockRepo := repository.NewMockAnnouncementRepository()
			h := NewHandler(service.NewAnnouncementService(mockRepo))
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupContext != nil {
				tt.setupContext(c)
			}

			h.AnnouncementsList(c, api.AnnouncementsListParams{
				IsActive: tt.isActive,
			})

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

// boolPtr は bool値のポインタを返すヘルパー関数
func boolPtr(b bool) *bool {
	return &b
}
