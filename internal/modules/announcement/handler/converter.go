package handler

import (
	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
)

func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:             announcement.ID,
		Title:          announcement.Title,
		AvailableFrom:  announcement.AvailableFrom,
		AvailableUntil: announcement.AvailableUntil,
		Url:            announcement.URL,
	}
}

func toDomainAnnouncementQuery(params api.AnnouncementsV1ListParams) domain.AnnouncementQuery {
	sortByDate := domain.SortDirectionAsc
	if params.SortByDate != nil {
		sortByDate = domain.SortDirection(*params.SortByDate)
	}

	filterIsActive := false
	if params.FilterIsActive != nil {
		filterIsActive = *params.FilterIsActive
	}

	return domain.AnnouncementQuery{
		FilterIsActive: filterIsActive,
		SortByDate:     sortByDate,
	}
}

func toDomainAnnouncementFromRequest(id string, req api.AnnouncementRequest) domain.Announcement {
	return domain.Announcement{
		ID:             id,
		Title:          req.Title,
		URL:            req.Url,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}
}
