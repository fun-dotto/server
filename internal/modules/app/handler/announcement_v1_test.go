package handler

import (
	"context"
	"testing"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/fun-dotto/server/internal/modules/app/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnouncementsV1List(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, resp api.AnnouncementsV1ListResponseObject, err error)
	}{
		{
			name: "正常にお知らせ一覧が取得できる",
			validate: func(t *testing.T, resp api.AnnouncementsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.AnnouncementsV1List200JSONResponse)
				require.True(t, ok, "レスポンスが200 JSONレスポンスではありません")
				assert.NotEmpty(t, result.Announcements, "アナウンスメントが空です")
			},
		},
		{
			name: "お知らせのフィールドが正しく返される",
			validate: func(t *testing.T, resp api.AnnouncementsV1ListResponseObject, err error) {
				require.NoError(t, err)
				result, ok := resp.(api.AnnouncementsV1List200JSONResponse)
				require.True(t, ok)
				assert.Len(t, result.Announcements, 1, "MockRepositoryは1件返すはずです")
				assert.Equal(t, "1", result.Announcements[0].Id)
				assert.Equal(t, "Announcement 1", result.Announcements[0].Title)
				assert.Equal(t, "https://example.com", result.Announcements[0].Url)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := repository.NewMockAnnouncementRepository()
			h := NewHandler(WithAnnouncementService(service.NewAnnouncementService(mockRepo)))

			resp, err := h.AnnouncementsV1List(context.Background(), api.AnnouncementsV1ListRequestObject{})

			if tt.validate != nil {
				tt.validate(t, resp, err)
			}
		})
	}
}
