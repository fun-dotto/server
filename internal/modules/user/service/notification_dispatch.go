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
		for _, t := range n.TargetUsers {
			userIDSet[t.UserID] = struct{}{}
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

	deliveries := make(map[string][]string, len(notifications))
	for _, n := range notifications {
		pendingUserIDs := make([]string, 0, len(n.TargetUsers))
		for _, t := range n.TargetUsers {
			pendingUserIDs = append(pendingUserIDs, t.UserID)
		}
		if len(pendingUserIDs) == 0 {
			continue
		}

		tokens, tokenUserIDs := collectTokens(pendingUserIDs, tokensByUser)
		if len(tokens) == 0 {
			deliveries[n.ID] = pendingUserIDs
			continue
		}

		successUserIDs, err := s.sendToTokens(ctx, n, tokens, tokenUserIDs)
		if err != nil {
			log.Printf("FCM send partially failed for notification %s (success=%d/%d users): %v", n.ID, len(successUserIDs), len(pendingUserIDs), err)
		}

		successSet := make(map[string]struct{}, len(successUserIDs))
		for _, uid := range successUserIDs {
			successSet[uid] = struct{}{}
		}
		delivered := make([]string, 0, len(pendingUserIDs))
		for _, uid := range pendingUserIDs {
			if _, ok := successSet[uid]; ok {
				delivered = append(delivered, uid)
				continue
			}
			if _, hasToken := tokensByUser[uid]; !hasToken {
				delivered = append(delivered, uid)
			}
		}
		if len(delivered) > 0 {
			deliveries[n.ID] = delivered
		}
	}

	if len(deliveries) == 0 {
		return []domain.Notification{}, nil
	}

	return s.repo.DispatchNotifications(ctx, deliveries)
}

func collectTokens(userIDs []string, tokensByUser map[string][]string) ([]string, []string) {
	seen := make(map[string]struct{})
	tokens := make([]string, 0)
	tokenUserIDs := make([]string, 0)
	for _, uid := range userIDs {
		for _, tk := range tokensByUser[uid] {
			if _, ok := seen[tk]; ok {
				continue
			}
			seen[tk] = struct{}{}
			tokens = append(tokens, tk)
			tokenUserIDs = append(tokenUserIDs, uid)
		}
	}
	return tokens, tokenUserIDs
}

const fcmMulticastBatchSize = 500

func (s *NotificationService) sendToTokens(ctx context.Context, n domain.Notification, tokens []string, tokenUserIDs []string) ([]string, error) {
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

	successUserSet := make(map[string]struct{})
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
			return collectSuccessUserIDs(tokenUserIDs, successUserSet), err
		}
		for i, r := range resp.Responses {
			uid := tokenUserIDs[start+i]
			if r.Error != nil {
				log.Printf("FCM delivery failed for notification %s token=%s: %v", n.ID, tokens[start+i], r.Error)
				continue
			}
			successUserSet[uid] = struct{}{}
		}
	}
	return collectSuccessUserIDs(tokenUserIDs, successUserSet), nil
}

func collectSuccessUserIDs(tokenUserIDs []string, successUserSet map[string]struct{}) []string {
	if len(successUserSet) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(successUserSet))
	result := make([]string, 0, len(successUserSet))
	for _, uid := range tokenUserIDs {
		if _, ok := successUserSet[uid]; !ok {
			continue
		}
		if _, dup := seen[uid]; dup {
			continue
		}
		seen[uid] = struct{}{}
		result = append(result, uid)
	}
	return result
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
