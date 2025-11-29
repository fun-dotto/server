package handler

import (
	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
)

func toApiAnnouncement(announcement domain.Announcement) api.Announcement {
	return api.Announcement{
		Id:       announcement.ID,
		Title:    announcement.Title,
		Date:     announcement.Date,
		Url:      announcement.URL,
		IsActive: announcement.IsActive,
	}
}

func toDomainAnnouncementQuery(params api.AnnouncementsListParams) domain.AnnouncementQuery {
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
