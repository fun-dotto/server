package service

import (
	"testing"
	"time"

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
			repo: repository.NewMockAnnouncementRepositoryWithError("getAnnouncements", assert.AnError),
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

func TestAnnouncementService_GetAnnouncement(t *testing.T) {
	tests := []struct {
		name     string
		repo     AnnouncementRepository
		id       string
		validate func(t *testing.T, announcement *domain.Announcement, err error)
	}{
		{
			name: "正常系: お知らせを取得できる",
			repo: repository.NewMockAnnouncementRepository(),
			id:   "1",
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.NoError(t, err)
				assert.Equal(t, "1", announcement.ID)
				assert.Equal(t, "Announcement 1", announcement.Title)
			},
		},
		{
			name: "異常系: 存在しないIDの場合エラーを返す",
			repo: repository.NewMockAnnouncementRepository(),
			id:   "nonexistent",
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.Error(t, err)
				assert.Nil(t, announcement)
			},
		},
		{
			name: "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo: repository.NewMockAnnouncementRepositoryWithError("getAnnouncement", assert.AnError),
			id:   "1",
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.Error(t, err)
				assert.Nil(t, announcement)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAnnouncementService(tt.repo)
			announcement, err := svc.GetAnnouncement(tt.id)
			tt.validate(t, announcement, err)
		})
	}
}

func TestAnnouncementService_CreateAnnouncement(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		repo     AnnouncementRepository
		req      domain.AnnouncementRequest
		validate func(t *testing.T, announcement *domain.Announcement, err error)
	}{
		{
			name: "正常系: お知らせを作成できる",
			repo: repository.NewMockAnnouncementRepository(),
			req: domain.AnnouncementRequest{
				Title:         "New Announcement",
				AvailableFrom: now,
				URL:           "https://example.com/new",
			},
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.NoError(t, err)
				assert.Equal(t, "New Announcement", announcement.Title)
				assert.Equal(t, "https://example.com/new", announcement.URL)
			},
		},
		{
			name: "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo: repository.NewMockAnnouncementRepositoryWithError("create", assert.AnError),
			req: domain.AnnouncementRequest{
				Title:         "New Announcement",
				AvailableFrom: now,
				URL:           "https://example.com/new",
			},
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.Error(t, err)
				assert.Nil(t, announcement)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAnnouncementService(tt.repo)
			announcement, err := svc.CreateAnnouncement(tt.req)
			tt.validate(t, announcement, err)
		})
	}
}

func TestAnnouncementService_UpdateAnnouncement(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		repo     AnnouncementRepository
		id       string
		req      domain.AnnouncementRequest
		validate func(t *testing.T, announcement *domain.Announcement, err error)
	}{
		{
			name: "正常系: お知らせを更新できる",
			repo: repository.NewMockAnnouncementRepository(),
			id:   "1",
			req: domain.AnnouncementRequest{
				Title:         "Updated Announcement",
				AvailableFrom: now,
				URL:           "https://example.com/updated",
			},
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.NoError(t, err)
				assert.Equal(t, "1", announcement.ID)
				assert.Equal(t, "Updated Announcement", announcement.Title)
				assert.Equal(t, "https://example.com/updated", announcement.URL)
			},
		},
		{
			name: "異常系: 存在しないIDの場合エラーを返す",
			repo: repository.NewMockAnnouncementRepository(),
			id:   "nonexistent",
			req: domain.AnnouncementRequest{
				Title:         "Updated",
				AvailableFrom: now,
				URL:           "https://example.com",
			},
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.Error(t, err)
				assert.Nil(t, announcement)
			},
		},
		{
			name: "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo: repository.NewMockAnnouncementRepositoryWithError("update", assert.AnError),
			id:   "1",
			req: domain.AnnouncementRequest{
				Title:         "Updated",
				AvailableFrom: now,
				URL:           "https://example.com",
			},
			validate: func(t *testing.T, announcement *domain.Announcement, err error) {
				require.Error(t, err)
				assert.Nil(t, announcement)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAnnouncementService(tt.repo)
			announcement, err := svc.UpdateAnnouncement(tt.id, tt.req)
			tt.validate(t, announcement, err)
		})
	}
}

func TestAnnouncementService_DeleteAnnouncement(t *testing.T) {
	tests := []struct {
		name     string
		repo     AnnouncementRepository
		id       string
		validate func(t *testing.T, err error)
	}{
		{
			name: "正常系: お知らせを削除できる",
			repo: repository.NewMockAnnouncementRepository(),
			id:   "1",
			validate: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "異常系: 存在しないIDの場合エラーを返す",
			repo: repository.NewMockAnnouncementRepository(),
			id:   "nonexistent",
			validate: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "異常系: リポジトリがエラーを返す場合エラーを返す",
			repo: repository.NewMockAnnouncementRepositoryWithError("delete", assert.AnError),
			id:   "1",
			validate: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAnnouncementService(tt.repo)
			err := svc.DeleteAnnouncement(tt.id)
			tt.validate(t, err)
		})
	}
}
