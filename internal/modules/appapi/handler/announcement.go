package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/shared/model"
)

// ListAnnouncements は公開中のおしらせの一覧を、公開開始日時の新しい順に返す。
func (h *Handler) ListAnnouncements(ctx context.Context, _ *connect.Request[appv1.ListAnnouncementsRequest]) (*connect.Response[appv1.ListAnnouncementsResponse], error) {
	now := h.now()
	var records []model.Announcement
	if err := h.db.WithContext(ctx).
		Where("available_from <= ?", now).
		Where("available_until IS NULL OR available_until > ?", now).
		Order("available_from DESC, id ASC").
		Find(&records).Error; err != nil {
		return nil, internalError(ctx, "ListAnnouncements", err)
	}
	return connect.NewResponse(&appv1.ListAnnouncementsResponse{
		Announcements: mapSlice(records, func(a model.Announcement) *appv1.Announcement {
			return &appv1.Announcement{
				Id:            a.ID.String(),
				Title:         a.Title,
				Url:           a.URL,
				AvailableFrom: toTimestamp(a.AvailableFrom),
			}
		}),
	}), nil
}
