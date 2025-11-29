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
		sortByDate   *api.SortDirection
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
				assert.Len(t, announcements, 2, "有効なお知らせは1件のはずです")
				assert.True(t, announcements[0].IsActive, "IsActiveがtrueではありません")
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
		{
			name:         "sortByDate=ascで日付昇順にソートされる",
			sortByDate:   sortDirPtr(api.Asc),
			isActive:     nil,
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(announcements), 2, "ソートテストには2件以上必要です")
				for i := 1; i < len(announcements); i++ {
					assert.True(t, !announcements[i].Date.Before(announcements[i-1].Date),
						"日付が昇順になっていません")
				}
			},
		},
		{
			name:         "sortByDate=descで日付降順にソートされる",
			sortByDate:   sortDirPtr(api.Desc),
			isActive:     nil,
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(announcements), 2, "ソートテストには2件以上必要です")
				for i := 1; i < len(announcements); i++ {
					assert.True(t, !announcements[i].Date.After(announcements[i-1].Date),
						"日付が降順になっていません")
				}
			},
		},
		{
			name:         "sortByDate未指定時はデフォルトで昇順",
			sortByDate:   nil,
			isActive:     nil,
			setupContext: func(c *gin.Context) {},
			wantCode:     http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var announcements []api.Announcement
				err := json.Unmarshal(w.Body.Bytes(), &announcements)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(announcements), 2, "ソートテストには2件以上必要です")
				for i := 1; i < len(announcements); i++ {
					assert.True(t, !announcements[i].Date.Before(announcements[i-1].Date),
						"デフォルトで昇順になっていません")
				}
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
				SortByDate:     tt.sortByDate,
				FilterIsActive: tt.isActive,
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

// sortDirPtr は SortDirection値のポインタを返すヘルパー関数
func sortDirPtr(s api.SortDirection) *api.SortDirection {
	return &s
}
