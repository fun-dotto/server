package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnnouncementHandler struct {
	c crud[model.Announcement]
}

func NewAnnouncementHandler(repo Repository[model.Announcement]) *AnnouncementHandler {
	return &AnnouncementHandler{c: crud[model.Announcement]{
		name: "announcement",
		repo: repo,
		keys: func(m model.Announcement) []keyField { return []keyField{{"id", m.ID.String()}} },
	}}
}

type announcementInput interface {
	GetTitle() string
	GetUrl() string
	GetAvailableFrom() *timestamppb.Timestamp
	GetAvailableUntil() *timestamppb.Timestamp
}

func announcementFromProto(in announcementInput, s fieldSet) (model.Announcement, error) {
	var v violations
	var m model.Announcement
	if s.has("title") {
		m.Title = requireString(&v, "title", in.GetTitle())
	}
	if s.has("url") {
		m.URL = requireString(&v, "url", in.GetUrl())
	}
	if s.has("available_from") {
		m.AvailableFrom = requireTimestamp(&v, "available_from", in.GetAvailableFrom())
	}
	if s.has("available_until") {
		m.AvailableUntil = optionalTimestamp(&v, "available_until", in.GetAvailableUntil())
	}
	return m, v.err()
}

func announcementToProto(m model.Announcement) *adminv1.Announcement {
	return &adminv1.Announcement{
		Id:             m.ID.String(),
		Title:          m.Title,
		Url:            m.URL,
		AvailableFrom:  toTimestamp(m.AvailableFrom),
		AvailableUntil: toOptionalTimestamp(m.AvailableUntil),
		CreatedAt:      toTimestamp(m.CreatedAt),
		UpdatedAt:      toTimestamp(m.UpdatedAt),
	}
}

func (h *AnnouncementHandler) ListAnnouncements(ctx context.Context, req *connect.Request[adminv1.ListAnnouncementsRequest]) (*connect.Response[adminv1.ListAnnouncementsResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListAnnouncementsResponse{NextPageToken: next}
	for _, m := range items {
		res.Announcements = append(res.Announcements, announcementToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *AnnouncementHandler) GetAnnouncement(ctx context.Context, req *connect.Request[adminv1.GetAnnouncementRequest]) (*connect.Response[adminv1.GetAnnouncementResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"id", id}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetAnnouncementResponse{Announcement: announcementToProto(m)}), nil
}

func (h *AnnouncementHandler) CreateAnnouncement(ctx context.Context, req *connect.Request[adminv1.CreateAnnouncementRequest]) (*connect.Response[adminv1.CreateAnnouncementResponse], error) {
	m, err := announcementFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateAnnouncementResponse{Announcement: announcementToProto(created)}), nil
}

func (h *AnnouncementHandler) UpdateAnnouncement(ctx context.Context, req *connect.Request[adminv1.UpdateAnnouncementRequest]) (*connect.Response[adminv1.UpdateAnnouncementResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "id")
	if err != nil {
		return nil, err
	}
	m, err := announcementFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"id", id}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateAnnouncementResponse{Announcement: announcementToProto(updated)}), nil
}

func (h *AnnouncementHandler) DeleteAnnouncement(ctx context.Context, req *connect.Request[adminv1.DeleteAnnouncementRequest]) (*connect.Response[adminv1.DeleteAnnouncementResponse], error) {
	id, err := parseIDField("id", req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"id", id}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteAnnouncementResponse{}), nil
}
