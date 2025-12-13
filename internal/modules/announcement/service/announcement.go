package service

import (
	"context"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type AnnouncementRepository interface {
	GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error)
}

type AnnouncementService struct {
	announcementRepository AnnouncementRepository
}

func NewAnnouncementService(announcementRepository AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepository: announcementRepository}
}

func (s *AnnouncementService) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	return s.announcementRepository.GetAnnouncements(ctx, query)
}
