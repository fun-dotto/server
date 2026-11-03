package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// AnnouncementHandler は AnnouncementService を実装する。
type AnnouncementHandler struct {
	appv1connect.UnimplementedAnnouncementServiceHandler
	service AnnouncementService
}

// NewAnnouncementHandler は AnnouncementHandler を返す。
func NewAnnouncementHandler(s AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{service: s}
}

// ListAnnouncements は公開中のおしらせの一覧を返す。
func (h *AnnouncementHandler) ListAnnouncements(ctx context.Context, req *connect.Request[appv1.ListAnnouncementsRequest]) (*connect.Response[appv1.ListAnnouncementsResponse], error) {
	announcements, err := h.service.GetAnnouncements()
	if err != nil {
		return nil, toConnectError(ctx, "ListAnnouncements", err)
	}
	page, next, err := paginate(req.Msg, announcements)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&appv1.ListAnnouncementsResponse{
		Announcements: mapSlice(page, toProtoAnnouncement),
		NextPageToken: next,
	}), nil
}

func toProtoAnnouncement(a domain.Announcement) *appv1.Announcement {
	return &appv1.Announcement{
		Id:            a.ID,
		Title:         a.Title,
		Url:           a.URL,
		AvailableFrom: toTimestamp(a.AvailableFrom),
	}
}
