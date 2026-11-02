package handler

import (
	"errors"

	"github.com/fun-dotto/server/internal/modules/app/service"
)

var (
	errAnnouncementServiceNotConfigured = errors.New("announcementService is not configured")
	errAcademicServiceNotConfigured     = errors.New("academicService is not configured")
	errUserServiceNotConfigured         = errors.New("userService is not configured")
)

type Handler struct {
	announcementService *service.AnnouncementService
	academicService     *service.AcademicService
	userService         *service.UserService
}

type Option func(*Handler)

func WithAnnouncementService(s *service.AnnouncementService) Option {
	return func(h *Handler) {
		h.announcementService = s
	}
}

func WithAcademicService(s *service.AcademicService) Option {
	return func(h *Handler) {
		h.academicService = s
	}
}

func WithUserService(s *service.UserService) Option {
	return func(h *Handler) {
		h.userService = s
	}
}

func NewHandler(opts ...Option) *Handler {
	h := &Handler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
