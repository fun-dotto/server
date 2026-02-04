package repository

import (
	"context"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type MockAnnouncementRepository struct {
	GetAnnouncementsFunc    func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error)
	GetAnnouncementByIDFunc func(ctx context.Context, id string) (domain.Announcement, error)
	CreateAnnouncementFunc  func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error)
	UpdateAnnouncementFunc  func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error)
	DeleteAnnouncementFunc  func(ctx context.Context, id string) error
}

func (m *MockAnnouncementRepository) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	if m.GetAnnouncementsFunc != nil {
		return m.GetAnnouncementsFunc(ctx, query)
	}
	return []domain.Announcement{}, nil
}

func (m *MockAnnouncementRepository) GetAnnouncementByID(ctx context.Context, id string) (domain.Announcement, error) {
	if m.GetAnnouncementByIDFunc != nil {
		return m.GetAnnouncementByIDFunc(ctx, id)
	}
	return domain.Announcement{}, nil
}

func (m *MockAnnouncementRepository) CreateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	if m.CreateAnnouncementFunc != nil {
		return m.CreateAnnouncementFunc(ctx, announcement)
	}
	return announcement, nil
}

func (m *MockAnnouncementRepository) UpdateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	if m.UpdateAnnouncementFunc != nil {
		return m.UpdateAnnouncementFunc(ctx, announcement)
	}
	return announcement, nil
}

func (m *MockAnnouncementRepository) DeleteAnnouncement(ctx context.Context, id string) error {
	if m.DeleteAnnouncementFunc != nil {
		return m.DeleteAnnouncementFunc(ctx, id)
	}
	return nil
}
