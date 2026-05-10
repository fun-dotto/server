package repository

import (
	"time"

	"github.com/fun-dotto/server/internal/modules/batch-jobs/domain"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

func parseDate(s string) time.Time {
	t, _ := time.Parse(dateLayout, s)
	return t
}

func subjectToDomain(m *model.Subject) domain.Subject {
	if m == nil {
		return domain.Subject{}
	}
	return domain.Subject{
		ID:   m.ID.String(),
		Name: m.Name,
	}
}

func roomToDomain(m *model.Room) domain.Room {
	if m == nil {
		return domain.Room{}
	}
	return domain.Room{
		ID:   m.ID.String(),
		Name: m.Name,
	}
}

func cancelledClassToDomain(m *model.CancelledClass) domain.CancelledClass {
	d := domain.CancelledClass{
		ID:     m.ID.String(),
		Date:   parseDate(m.Date),
		Period: m.Period,
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(m.Subject)
	}
	return d
}

func makeupClassToDomain(m *model.MakeupClass) domain.MakeupClass {
	d := domain.MakeupClass{
		ID:     m.ID.String(),
		Date:   parseDate(m.Date),
		Period: m.Period,
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(m.Subject)
	}
	return d
}

func roomChangeToDomain(m *model.RoomChange) domain.RoomChange {
	d := domain.RoomChange{
		ID:     m.ID.String(),
		Date:   parseDate(m.Date),
		Period: m.Period,
	}
	if m.Subject != nil {
		d.Subject = subjectToDomain(m.Subject)
	}
	if m.NewRoom != nil {
		d.NewRoom = roomToDomain(m.NewRoom)
	}
	return d
}

func fcmTokenToDomain(m *model.FCMToken) domain.FCMToken {
	return domain.FCMToken{
		Token:     m.Token,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func notificationToDomain(m *model.Notification, targets []domain.NotificationTargetUser) domain.Notification {
	return domain.Notification{
		ID:                   m.ID.String(),
		Title:                m.Title,
		Body:                 m.Body,
		ImageURL:             m.ImageURL,
		AnalyticsLabel:       m.AnalyticsLabel,
		APNsBadge:            m.APNsBadge,
		APNsSound:            m.APNsSound,
		APNsContentAvailable: m.APNsContentAvailable,
		AndroidChannelID:     m.AndroidChannelID,
		AndroidPriority:      m.AndroidPriority,
		AndroidTTLSeconds:    m.AndroidTTLSeconds,
		WebpushLink:          m.WebpushLink,
		URL:                  m.URL,
		NotifyAfter:          m.NotifyAfter,
		NotifyBefore:         m.NotifyBefore,
		TargetUsers:          targets,
	}
}

func notificationFromDomain(n domain.Notification) (model.Notification, error) {
	id, err := uuid.Parse(n.ID)
	if err != nil {
		return model.Notification{}, err
	}
	m := model.Notification{
		Title:                n.Title,
		Body:                 n.Body,
		ImageURL:             n.ImageURL,
		AnalyticsLabel:       n.AnalyticsLabel,
		APNsBadge:            n.APNsBadge,
		APNsSound:            n.APNsSound,
		APNsContentAvailable: n.APNsContentAvailable,
		AndroidChannelID:     n.AndroidChannelID,
		AndroidPriority:      n.AndroidPriority,
		AndroidTTLSeconds:    n.AndroidTTLSeconds,
		WebpushLink:          n.WebpushLink,
		URL:                  n.URL,
		NotifyAfter:          n.NotifyAfter,
		NotifyBefore:         n.NotifyBefore,
	}
	m.ID = id
	return m, nil
}
