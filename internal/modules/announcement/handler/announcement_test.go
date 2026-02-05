package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/fun-dotto/announcement-api/internal/repository"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsV1List(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)

	tests := []struct {
		name      string
		setupMock func() *repository.MockAnnouncementRepository
		params    api.AnnouncementsV1ListParams
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせ一覧が取得できる",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
						return []domain.Announcement{
							{ID: "1", Title: "お知らせ1", AvailableFrom: now, URL: "https://example.com/1"},
							{ID: "2", Title: "お知らせ2", AvailableFrom: yesterday, URL: "https://example.com/2"},
							{ID: "3", Title: "お知らせ3", AvailableFrom: twoDaysAgo, URL: "https://example.com/3"},
						}, nil
					},
				}
			},
			params:   api.AnnouncementsV1ListParams{},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.AnnouncementsV1List200JSONResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "failed to unmarshal response body")
				assert.Len(t, response.Announcements, 3)
			},
		},
		{
			name: "空の結果を正常に返せる",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
						return []domain.Announcement{}, nil
					},
				}
			},
			params:   api.AnnouncementsV1ListParams{},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.AnnouncementsV1List200JSONResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Empty(t, response.Announcements)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			h := NewHandler(service.NewAnnouncementService(mockRepo))

			request := api.AnnouncementsV1ListRequestObject{Params: tt.params}
			response, err := h.AnnouncementsV1List(context.Background(), request)

			assert.NoError(t, err)
			assert.NotNil(t, response)

			w := httptest.NewRecorder()
			err = response.VisitAnnouncementsV1ListResponse(w)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
