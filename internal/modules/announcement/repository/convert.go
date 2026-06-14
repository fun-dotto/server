package repository

import (
	"github.com/fun-dotto/server/internal/modules/announcement/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
)

// shared/model は ID を uuid.UUID で保持する一方、announcement の domain 層は
// 文字列 ID を扱う。境界変換をこのファイルに集約する。

func parseUUIDOrNil(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

func toDomainAnnouncement(m model.Announcement) domain.Announcement {
	return domain.Announcement{
		ID:             m.ID.String(),
		Title:          m.Title,
		URL:            m.URL,
		AvailableFrom:  m.AvailableFrom,
		AvailableUntil: m.AvailableUntil,
	}
}

func fromDomainAnnouncement(a domain.Announcement) model.Announcement {
	return model.Announcement{
		Common:         model.Common{ID: parseUUIDOrNil(a.ID)},
		Title:          a.Title,
		URL:            a.URL,
		AvailableFrom:  a.AvailableFrom,
		AvailableUntil: a.AvailableUntil,
	}
}
