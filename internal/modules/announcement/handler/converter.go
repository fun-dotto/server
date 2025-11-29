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
	return domain.AnnouncementQuery{
		IsActive: params.IsActive,
	}
}
