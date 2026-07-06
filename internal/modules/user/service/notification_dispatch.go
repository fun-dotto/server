package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/fun-dotto/server/internal/modules/user/domain"
)

// DispatchNotifications は指定 ID の通知を対象ユーザーへ FCM 送信し、notified_at を更新する (管理 API 用)。
func (s *NotificationService) DispatchNotifications(ctx context.Context, ids []string) ([]domain.Notification, error) {
	notifications, err := s.repo.GetNotificationsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(notifications) == 0 {
		return []domain.Notification{}, nil
	}

	tokensByUser, err := s.listTokensByUser(ctx, notifications)
	if err != nil {
		return nil, err
	}

	deliveries := make(map[string][]string, len(notifications))
	for _, n := range notifications {
		pendingUserIDs := targetUserIDs(n)
		if len(pendingUserIDs) == 0 {
			continue
		}

		tokens, tokenUserIDs := collectTokens(pendingUserIDs, tokensByUser)
		if len(tokens) == 0 {
			deliveries[n.ID] = pendingUserIDs
			continue
		}

		successUserIDs, _, err := s.sendToTokens(ctx, n, tokens, tokenUserIDs)
		if err != nil {
			log.Printf("FCM send partially failed for notification %s (success=%d/%d users): %v", n.ID, len(successUserIDs), len(pendingUserIDs), err)
		}

		delivered := deliveredUserIDs(pendingUserIDs, successUserIDs, tokensByUser)
		if len(delivered) > 0 {
			deliveries[n.ID] = delivered
		}
	}

	if len(deliveries) == 0 {
		return []domain.Notification{}, nil
	}

	return s.repo.DispatchNotifications(ctx, deliveries)
}

type DispatchSummary struct {
	Pending      int
	Dispatched   int
	NoTokenSkip  int
	FailedSend   int
	TotalFCMSent int
	DryRun       bool
}

// DispatchPendingNotifications は通知期間内で未通知ユーザーが残っている通知を FCM 送信する (バッチジョブ用)。
func (s *NotificationService) DispatchPendingNotifications(ctx context.Context, dryRun bool) (DispatchSummary, error) {
	summary := DispatchSummary{DryRun: dryRun}

	now := time.Now()
	notifications, err := s.repo.ListPendingNotifications(ctx, now)
	if err != nil {
		return summary, fmt.Errorf("list pending notifications: %w", err)
	}
	summary.Pending = len(notifications)
	if len(notifications) == 0 {
		return summary, nil
	}

	tokensByUser, err := s.listTokensByUser(ctx, notifications)
	if err != nil {
		return summary, err
	}

	for _, n := range notifications {
		pendingUserIDs := targetUserIDs(n)
		if len(pendingUserIDs) == 0 {
			continue
		}

		tokens, tokenUserIDs := collectTokens(pendingUserIDs, tokensByUser)
		if len(tokens) == 0 {
			// FCM トークン未登録のユーザーは送りようが無いので、再ディスパッチを止めるため通知済み扱いにする。
			log.Printf("skip: no fcm_tokens for notification %s (targets=%d) — marking %d user(s) as delivered to avoid re-dispatch", n.ID, len(pendingUserIDs), len(pendingUserIDs))
			summary.NoTokenSkip++
			if dryRun {
				continue
			}
			if err := s.repo.MarkUsersAsNotified(ctx, map[string][]string{n.ID: pendingUserIDs}); err != nil {
				return summary, fmt.Errorf("mark users as notified %s: %w", n.ID, err)
			}
			continue
		}

		if dryRun {
			log.Printf("dry-run: would send notification %s to %d token(s) (targets=%d) title=%q", n.ID, len(tokens), len(pendingUserIDs), n.Title)
			summary.TotalFCMSent += len(tokens)
			summary.Dispatched++
			continue
		}

		successUserIDs, sentTokens, err := s.sendToTokens(ctx, n, tokens, tokenUserIDs)
		summary.TotalFCMSent += sentTokens
		if err != nil {
			log.Printf("FCM send partially failed for notification %s (success=%d/%d users, tokens=%d/%d): %v", n.ID, len(successUserIDs), len(pendingUserIDs), sentTokens, len(tokens), err)
		}

		delivered := deliveredUserIDs(pendingUserIDs, successUserIDs, tokensByUser)
		if len(delivered) == 0 {
			summary.FailedSend++
			continue
		}
		if err := s.repo.MarkUsersAsNotified(ctx, map[string][]string{n.ID: delivered}); err != nil {
			return summary, fmt.Errorf("mark users as notified %s: %w", n.ID, err)
		}
		summary.Dispatched++
	}

	return summary, nil
}

// listTokensByUser は通知の対象ユーザー全員分の FCM トークンをユーザー ID 毎にまとめて取得する。
func (s *NotificationService) listTokensByUser(ctx context.Context, notifications []domain.Notification) (map[string][]string, error) {
	userIDSet := make(map[string]struct{})
	for _, n := range notifications {
		for _, t := range n.TargetUsers {
			userIDSet[t.UserID] = struct{}{}
		}
	}

	tokensByUser := make(map[string][]string)
	if len(userIDSet) == 0 {
		return tokensByUser, nil
	}

	userIDs := make([]string, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}
	tokens, err := s.fcmTokenRepo.ListFCMTokens(ctx, domain.FCMTokenListFilter{UserIDs: userIDs})
	if err != nil {
		return nil, fmt.Errorf("list fcm_tokens: %w", err)
	}
	for _, t := range tokens {
		tokensByUser[t.UserID] = append(tokensByUser[t.UserID], t.Token)
	}
	return tokensByUser, nil
}

func targetUserIDs(n domain.Notification) []string {
	userIDs := make([]string, 0, len(n.TargetUsers))
	for _, t := range n.TargetUsers {
		userIDs = append(userIDs, t.UserID)
	}
	return userIDs
}

// deliveredUserIDs は送信成功したユーザーに加え、トークンが1件も登録されていない
// (再送しても届かない) ユーザーを配信扱いとして返す。
func deliveredUserIDs(pendingUserIDs, successUserIDs []string, tokensByUser map[string][]string) []string {
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
	return delivered
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

func (s *NotificationService) sendToTokens(ctx context.Context, n domain.Notification, tokens []string, tokenUserIDs []string) ([]string, int, error) {
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
	sentTokens := 0
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
			return collectSuccessUserIDs(tokenUserIDs, successUserSet), sentTokens, err
		}
		for i, r := range resp.Responses {
			uid := tokenUserIDs[start+i]
			if r.Error != nil {
				log.Printf("FCM delivery failed for notification %s token=%s: %v", n.ID, redactToken(tokens[start+i]), r.Error)
				continue
			}
			sentTokens++
			successUserSet[uid] = struct{}{}
		}
	}
	return collectSuccessUserIDs(tokenUserIDs, successUserSet), sentTokens, nil
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

// redactToken は FCM トークンをログ出力する際に、先頭と末尾の数文字だけ残して個人識別につながる情報量を削る。
func redactToken(token string) string {
	const keep = 4
	if len(token) <= keep*2 {
		return "***"
	}
	return token[:keep] + "..." + token[len(token)-keep:]
}
