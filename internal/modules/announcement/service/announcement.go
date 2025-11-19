package service

import "github.com/fun-dotto/announcement-api/internal/domain"

type AnnouncementRepository interface {
	GetAnnouncements(isActive *bool) ([]domain.Announcement, error)
}

type AnnouncementService struct {
	announcementRepository AnnouncementRepository
}

func NewAnnouncementService(announcementRepository AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepository: announcementRepository}
}

func (s *AnnouncementService) GetAnnouncements(isActive *bool) ([]domain.Announcement, error) {
	return s.announcementRepository.GetAnnouncements(isActive)
}
