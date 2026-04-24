package service

import (
	"context"
	"log"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/fun-dotto/user-api/internal/domain"
)

func (s *NotificationService) DispatchNotifications(ctx context.Context, ids []string) ([]domain.Notification, error) {
	notifications, err := s.repo.GetNotificationsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(notifications) == 0 {
		return []domain.Notification{}, nil
	}

	userIDSet := make(map[string]struct{})
	for _, n := range notifications {
		for _, uid := range n.TargetUserIDs {
			userIDSet[uid] = struct{}{}
		}
	}

	tokensByUser := make(map[string][]string)
	if len(userIDSet) > 0 {
		userIDs := make([]string, 0, len(userIDSet))
		for uid := range userIDSet {
			userIDs = append(userIDs, uid)
		}
		tokens, err := s.fcmTokenRepo.ListFCMTokens(ctx, domain.FCMTokenListFilter{UserIDs: userIDs})
		if err != nil {
			return nil, err
		}
		for _, t := range tokens {
			tokensByUser[t.UserID] = append(tokensByUser[t.UserID], t.Token)
		}
	}

	successIDs := make([]string, 0, len(notifications))
	for _, n := range notifications {
		tokens := collectTokens(n.TargetUserIDs, tokensByUser)
		if len(tokens) == 0 {
			successIDs = append(successIDs, n.ID)
			continue
		}

		sent, err := s.sendToTokens(ctx, n, tokens)
		if err != nil {
			log.Printf("FCM send failed for notification %s: %v", n.ID, err)
			continue
		}
		if sent > 0 {
			successIDs = append(successIDs, n.ID)
		}
	}

	if len(successIDs) == 0 {
		return []domain.Notification{}, nil
	}

	return s.repo.DispatchNotifications(ctx, successIDs)
}

func collectTokens(userIDs []string, tokensByUser map[string][]string) []string {
	seen := make(map[string]struct{})
	tokens := make([]string, 0)
	for _, uid := range userIDs {
		for _, tk := range tokensByUser[uid] {
			if _, ok := seen[tk]; ok {
				continue
			}
			seen[tk] = struct{}{}
			tokens = append(tokens, tk)
		}
	}
	return tokens
}

const fcmMulticastBatchSize = 500

func (s *NotificationService) sendToTokens(ctx context.Context, n domain.Notification, tokens []string) (int, error) {
	data := map[string]string{"notification_id": n.ID}
	if n.URL != nil {
		data["url"] = *n.URL
	}

	notification := &messaging.Notification{
		Title: n.Title,
		Body:  n.Body,
	}
	if n.ImageURL != nil {
		notification.ImageURL = *n.ImageURL
	}

	var fcmOptions *messaging.FCMOptions
	if n.AnalyticsLabel != nil {
		fcmOptions = &messaging.FCMOptions{AnalyticsLabel: *n.AnalyticsLabel}
	}

	androidConfig := buildAndroidConfig(n)
	apnsConfig := buildAPNSConfig(n)
	webpushConfig := buildWebpushConfig(n)

	totalSuccess := 0
	for start := 0; start < len(tokens); start += fcmMulticastBatchSize {
		end := min(start+fcmMulticastBatchSize, len(tokens))
		msg := &messaging.MulticastMessage{
			Tokens:       tokens[start:end],
			Notification: notification,
			Data:         data,
			Android:      androidConfig,
			APNS:         apnsConfig,
			Webpush:      webpushConfig,
			FCMOptions:   fcmOptions,
		}
		resp, err := s.messagingClient.SendEachForMulticast(ctx, msg)
		if err != nil {
			return totalSuccess, err
		}
		totalSuccess += resp.SuccessCount
		if resp.FailureCount > 0 {
			for i, r := range resp.Responses {
				if r.Error != nil {
					log.Printf("FCM delivery failed for notification %s token=%s: %v", n.ID, tokens[start+i], r.Error)
				}
			}
		}
	}
	return totalSuccess, nil
}

func buildAndroidConfig(n domain.Notification) *messaging.AndroidConfig {
	if n.AndroidChannelID == nil && n.AndroidPriority == nil && n.AndroidTTLSeconds == nil {
		return nil
	}
	cfg := &messaging.AndroidConfig{}
	if n.AndroidPriority != nil {
		cfg.Priority = *n.AndroidPriority
	}
	if n.AndroidTTLSeconds != nil {
		ttl := time.Duration(*n.AndroidTTLSeconds) * time.Second
		cfg.TTL = &ttl
	}
	if n.AndroidChannelID != nil {
		cfg.Notification = &messaging.AndroidNotification{ChannelID: *n.AndroidChannelID}
	}
	return cfg
}

func buildAPNSConfig(n domain.Notification) *messaging.APNSConfig {
	if n.APNsBadge == nil && n.APNsSound == nil && n.APNsContentAvailable == nil {
		return nil
	}
	aps := &messaging.Aps{}
	if n.APNsBadge != nil {
		aps.Badge = n.APNsBadge
	}
	if n.APNsSound != nil {
		aps.Sound = *n.APNsSound
	}
	if n.APNsContentAvailable != nil {
		aps.ContentAvailable = *n.APNsContentAvailable
	}
	return &messaging.APNSConfig{Payload: &messaging.APNSPayload{Aps: aps}}
}

func buildWebpushConfig(n domain.Notification) *messaging.WebpushConfig {
	if n.WebpushLink == nil {
		return nil
	}
	return &messaging.WebpushConfig{
		FCMOptions: &messaging.WebpushFCMOptions{Link: *n.WebpushLink},
	}
}
