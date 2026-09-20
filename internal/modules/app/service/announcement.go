package service

import "github.com/fun-dotto/server/internal/modules/app/domain"

type AnnouncementRepository interface {
	GetAnnouncements(query domain.AnnouncementQuery) ([]domain.Announcement, error)
	GetAnnouncement(id string) (*domain.Announcement, error)
	CreateAnnouncement(req domain.AnnouncementRequest) (*domain.Announcement, error)
	UpdateAnnouncement(id string, req domain.AnnouncementRequest) (*domain.Announcement, error)
	DeleteAnnouncement(id string) error
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

func (s *AnnouncementService) GetAnnouncement(id string) (*domain.Announcement, error) {
	return s.announcementRepository.GetAnnouncement(id)
}

func (s *AnnouncementService) CreateAnnouncement(req domain.AnnouncementRequest) (*domain.Announcement, error) {
	return s.announcementRepository.CreateAnnouncement(req)
}

func (s *AnnouncementService) UpdateAnnouncement(id string, req domain.AnnouncementRequest) (*domain.Announcement, error) {
	return s.announcementRepository.UpdateAnnouncement(id, req)
}

func (s *AnnouncementService) DeleteAnnouncement(id string) error {
	return s.announcementRepository.DeleteAnnouncement(id)
}
