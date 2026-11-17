package service

import "github.com/fun-dotto/server/internal/modules/app/domain"

type AnnouncementRepository interface {
	GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error)
}

type AnnouncementService struct {
	announcementRepository AnnouncementRepository
}

func NewAnnouncementService(announcementRepository AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepository: announcementRepository}
}

func (s *AnnouncementService) GetAnnouncements() ([]domain.Announcement, error) {
	sortByDate := domain.SortDirectionDesc
	filterIsActive := true
	query := domain.AnnouncementQuery{
		SortByDate:     &sortByDate,
		FilterIsActive: &filterIsActive,
	}
	return s.announcementRepository.GetAnnouncements(query)
}
