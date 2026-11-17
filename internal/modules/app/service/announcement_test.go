package service

import (
	"testing"

	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnnouncementService_GetAnnouncements(t *testing.T) {
	tests := []struct {
		name     string
		repo     AnnouncementRepository
		validate func(t *testing.T, announcements []domain.Announcement, err error)
	}{
		{
			name: "正常系: お知らせ一覧を取得できる",
			repo: repository.NewMockAnnouncementRepository(),
			validate: func(t *testing.T, announcements []domain.Announcement, err error) {
				require.NoError(t, err)
				assert.Len(t, announcements, 1)
				assert.Equal(t, "1", announcements[0].ID)
				assert.Equal(t, "Announcement 1", announcements[0].Title)
			},
		},
		{
			name: "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo: repository.NewMockAnnouncementRepositoryWithError(assert.AnError),
			validate: func(t *testing.T, announcements []domain.Announcement, err error) {
				require.Error(t, err)
				assert.Nil(t, announcements)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAnnouncementService(tt.repo)
			announcements, err := svc.GetAnnouncements()
			tt.validate(t, announcements, err)
		})
	}
}
