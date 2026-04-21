package service

import (
	"context"
	"log"

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

	totalSuccess := 0
	for start := 0; start < len(tokens); start += fcmMulticastBatchSize {
		end := min(start+fcmMulticastBatchSize, len(tokens))
		msg := &messaging.MulticastMessage{
			Tokens: tokens[start:end],
			Notification: &messaging.Notification{
				Title: n.Title,
				Body:  n.Message,
			},
			Data: data,
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
