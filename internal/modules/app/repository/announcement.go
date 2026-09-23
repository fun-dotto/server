package repository

import (
	"context"
	"fmt"

	announcement_api "github.com/fun-dotto/server/gen/announcement"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/external"
)

// AnnouncementRepository は外部APIからお知らせを取得する
type AnnouncementRepository struct {
	client *announcement_api.ClientWithResponses
}

func NewAnnouncementRepository(client *announcement_api.ClientWithResponses) *AnnouncementRepository {
	return &AnnouncementRepository{client: client}
}

func (r *AnnouncementRepository) GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	params := announcement_api.AnnouncementsV1ListParams{
		SortByDate: func() *announcement_api.FoundationV1SortDirection {
			if query.SortByDate == nil {
				return nil
			}
			sortDirection := announcement_api.FoundationV1SortDirection(*query.SortByDate)
			return &sortDirection
		}(),
		FilterIsActive: query.FilterIsActive,
	}
	// 外部APIからデータ取得
	response, err := r.client.AnnouncementsV1ListWithResponse(context.Background(), &params)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get announcements: status %d", response.StatusCode())
	}

	// 外部API形式 → ドメイン形式に変換
	announcements := response.JSON200.Announcements
	result := make([]domain.Announcement, len(announcements))
	for i, a := range announcements {
		result[i] = external.ToDomainAnnouncement(a)
	}

	return result, nil
}

func (r *AnnouncementRepository) GetAnnouncement(id string) (*domain.Announcement, error) {
	response, err := r.client.AnnouncementsV1DetailWithResponse(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get announcement: status %d", response.StatusCode())
	}

	a := external.ToDomainAnnouncement(response.JSON200.Announcement)
	return &a, nil
}

func (r *AnnouncementRepository) CreateAnnouncement(req domain.AnnouncementRequest) (*domain.Announcement, error) {
	body := external.ToExternalAnnouncementRequest(req)
	response, err := r.client.AnnouncementsV1CreateWithResponse(context.Background(), body)
	if err != nil {
		return nil, err
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create announcement: status %d", response.StatusCode())
	}

	a := external.ToDomainAnnouncement(response.JSON201.Announcement)
	return &a, nil
}

func (r *AnnouncementRepository) UpdateAnnouncement(id string, req domain.AnnouncementRequest) (*domain.Announcement, error) {
	body := external.ToExternalAnnouncementRequest(req)
	response, err := r.client.AnnouncementsV1UpdateWithResponse(context.Background(), id, body)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update announcement: status %d", response.StatusCode())
	}

	a := external.ToDomainAnnouncement(response.JSON200.Announcement)
	return &a, nil
}

func (r *AnnouncementRepository) DeleteAnnouncement(id string) error {
	response, err := r.client.AnnouncementsV1DeleteWithResponse(context.Background(), id)
	if err != nil {
		return err
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete announcement: status %d", response.StatusCode())
	}

	return nil
}
