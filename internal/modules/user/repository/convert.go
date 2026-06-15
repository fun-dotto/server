package repository

import (
	"github.com/fun-dotto/server/internal/modules/user/domain"
	"github.com/fun-dotto/server/internal/shared/model"
)

func userToDomain(m model.User) domain.User {
	return domain.User{
		ID:     m.ID,
		Email:  m.Email,
		Grade:  m.Grade,
		Course: m.Course,
		Class:  m.Class,
	}
}

func userFromDomain(u domain.User) model.User {
	return model.User{
		ID:     u.ID,
		Email:  u.Email,
		Grade:  u.Grade,
		Course: u.Course,
		Class:  u.Class,
	}
}

func fcmTokenToDomain(m model.FCMToken) domain.FCMToken {
	return domain.FCMToken{
		Token:     m.Token,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func fcmTokenFromDomain(t domain.FCMToken) model.FCMToken {
	return model.FCMToken{
		Token:     t.Token,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func notificationToDomain(m model.Notification, targets []domain.NotificationTargetUser) domain.Notification {
	return domain.Notification{
		ID:                   m.ID,
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

func notificationFromDomain(n domain.Notification) model.Notification {
	return model.Notification{
		ID:                   n.ID,
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
}
