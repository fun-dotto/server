package service

import "github.com/fun-dotto/announcement-api/internal/domain"

type AnouncementRepository interface {
	GetAnnouncements() ([]domain.Announcement, error)
}

type AnnouncementService struct {
	announcementRepository AnouncementRepository
}

func NewAnnouncementService(announcementRepository AnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepository: announcementRepository}
}

func (s *AnnouncementService) GetAnnouncements() ([]domain.Announcement, error) {
	return s.announcementRepository.GetAnnouncements()
}
