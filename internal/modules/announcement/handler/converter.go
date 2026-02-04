package handler

import (
	"time"

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
		Date:           announcement.Date,
		IsActive:       announcement.IsActive,
	}
}

// v0廃止まで残す
func toDomainAnnouncementQuery(params api.AnnouncementsV0ListParams) domain.AnnouncementQuery {
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

func toDomainAnnouncementQueryV1(params api.AnnouncementsV1ListParams) domain.AnnouncementQuery {
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

func toDomainAnnouncementFromRequest(id string, req api.AnnouncementRequest, now time.Time) domain.Announcement {
	// v0との後方互換性: IsActiveはAvailableUntilがnilまたは未来の場合にtrue
	isActive := req.AvailableUntil == nil || req.AvailableUntil.After(now)

	return domain.Announcement{
		ID:             id,
		Title:          req.Title,
		URL:            req.Url,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
		Date:           req.AvailableFrom, // v0との後方互換性
		IsActive:       isActive,
	}
}
