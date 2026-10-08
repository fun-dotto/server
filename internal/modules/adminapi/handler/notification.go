package handler

import (
	"context"
	"uuid"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationHandler struct {
	c crud[model.Notification]
}

func NewNotificationHandler(repo Repository[model.Notification]) *NotificationHandler {
	return &NotificationHandler{c: crud[model.Notification]{
		name: "notification",
		repo: repo,
		keys: func(m model.Notification) []keyField { return []keyField{{"id", m.ID}} },
	}}
}

var notificationAndroidPriorityEnum = enumMap[adminv1.NotificationAndroidPriority]{
	adminv1.NotificationAndroidPriority_NOTIFICATION_ANDROID_PRIORITY_NORMAL: "normal",
	adminv1.NotificationAndroidPriority_NOTIFICATION_ANDROID_PRIORITY_HIGH:   "high",
}

// notificationColumns は GORM の命名規則で列名が proto のフィールド名と一致しないもの。
var notificationColumns = map[string][]string{
	"apns_badge":             {"APNsBadge"},
	"apns_sound":             {"APNsSound"},
	"apns_content_available": {"APNsContentAvailable"},
	"android_ttl_seconds":    {"AndroidTTLSeconds"},
}

type notificationInput interface {
	proto.Message
	GetTitle() string
	GetBody() string
	GetImageUrl() string
	GetAnalyticsLabel() string
	GetApnsBadge() int32
	GetApnsSound() string
	GetApnsContentAvailable() bool
	GetAndroidChannelId() string
	GetAndroidPriority() adminv1.NotificationAndroidPriority
	GetAndroidTtlSeconds() int32
	GetWebpushLink() string
	GetUrl() string
	GetNotifyAfter() *timestamppb.Timestamp
	GetNotifyBefore() *timestamppb.Timestamp
}

func notificationFromProto(in notificationInput, s fieldSet) (model.Notification, error) {
	var v violations
	var m model.Notification
	if s.has("title") {
		m.Title = requireString(&v, "title", in.GetTitle())
	}
	if s.has("body") {
		m.Body = requireString(&v, "body", in.GetBody())
	}
	if s.has("image_url") {
		m.ImageURL = optional(in, "image_url", in.GetImageUrl())
	}
	if s.has("analytics_label") {
		m.AnalyticsLabel = optional(in, "analytics_label", in.GetAnalyticsLabel())
	}
	if s.has("apns_badge") {
		m.APNsBadge = intPtr(optional(in, "apns_badge", in.GetApnsBadge()))
	}
	if s.has("apns_sound") {
		m.APNsSound = optional(in, "apns_sound", in.GetApnsSound())
	}
	if s.has("apns_content_available") {
		m.APNsContentAvailable = optional(in, "apns_content_available", in.GetApnsContentAvailable())
	}
	if s.has("android_channel_id") {
		m.AndroidChannelID = optional(in, "android_channel_id", in.GetAndroidChannelId())
	}
	if s.has("android_priority") {
		m.AndroidPriority = notificationAndroidPriorityEnum.toOptionalDB(&v, "android_priority",
			optional(in, "android_priority", in.GetAndroidPriority()))
	}
	if s.has("android_ttl_seconds") {
		m.AndroidTTLSeconds = intPtr(optional(in, "android_ttl_seconds", in.GetAndroidTtlSeconds()))
	}
	if s.has("webpush_link") {
		m.WebpushLink = optional(in, "webpush_link", in.GetWebpushLink())
	}
	if s.has("url") {
		m.URL = optional(in, "url", in.GetUrl())
	}
	if s.has("notify_after") {
		m.NotifyAfter = requireTimestamp(&v, "notify_after", in.GetNotifyAfter())
	}
	if s.has("notify_before") {
		m.NotifyBefore = requireTimestamp(&v, "notify_before", in.GetNotifyBefore())
	}
	return m, v.err()
}

func notificationToProto(m model.Notification) *adminv1.Notification {
	return &adminv1.Notification{
		Id:                   m.ID,
		Title:                m.Title,
		Body:                 m.Body,
		ImageUrl:             m.ImageURL,
		AnalyticsLabel:       m.AnalyticsLabel,
		ApnsBadge:            int32Ptr(m.APNsBadge),
		ApnsSound:            m.APNsSound,
		ApnsContentAvailable: m.APNsContentAvailable,
		AndroidChannelId:     m.AndroidChannelID,
		AndroidPriority:      notificationAndroidPriorityEnum.fromOptionalDB(m.AndroidPriority),
		AndroidTtlSeconds:    int32Ptr(m.AndroidTTLSeconds),
		WebpushLink:          m.WebpushLink,
		Url:                  m.URL,
		NotifyAfter:          toTimestamp(m.NotifyAfter),
		NotifyBefore:         toTimestamp(m.NotifyBefore),
		CreatedAt:            toTimestamp(m.CreatedAt),
		UpdatedAt:            toTimestamp(m.UpdatedAt),
	}
}

func (h *NotificationHandler) ListNotifications(ctx context.Context, req *connect.Request[adminv1.ListNotificationsRequest]) (*connect.Response[adminv1.ListNotificationsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListNotificationsResponse{NextPageToken: next}
	for _, m := range items {
		res.Notifications = append(res.Notifications, notificationToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *NotificationHandler) GetNotification(ctx context.Context, req *connect.Request[adminv1.GetNotificationRequest]) (*connect.Response[adminv1.GetNotificationResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetNotificationResponse{Notification: notificationToProto(m)}), nil
}

func (h *NotificationHandler) CreateNotification(ctx context.Context, req *connect.Request[adminv1.CreateNotificationRequest]) (*connect.Response[adminv1.CreateNotificationResponse], error) {
	m, err := notificationFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	m.ID = uuid.New().String()
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateNotificationResponse{Notification: notificationToProto(created)}), nil
}

func (h *NotificationHandler) UpdateNotification(ctx context.Context, req *connect.Request[adminv1.UpdateNotificationRequest]) (*connect.Response[adminv1.UpdateNotificationResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := notificationFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(notificationColumns))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateNotificationResponse{Notification: notificationToProto(updated)}), nil
}

func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *connect.Request[adminv1.DeleteNotificationRequest]) (*connect.Response[adminv1.DeleteNotificationResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteNotificationResponse{}), nil
}
