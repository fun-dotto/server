package handler

import (
	academic_api "github.com/fun-dotto/server/gen/academic"
	api "github.com/fun-dotto/server/gen/admin"
	announcement_api "github.com/fun-dotto/server/gen/announcement"
	funch_api "github.com/fun-dotto/server/gen/funch"
	user_api "github.com/fun-dotto/server/gen/user"
)

type Handler struct {
	academicClient     *academic_api.ClientWithResponses
	announcementClient *announcement_api.ClientWithResponses
	funchClient        *funch_api.ClientWithResponses
	userClient         *user_api.ClientWithResponses
}

func NewHandler(
	academicClient *academic_api.ClientWithResponses,
	announcementClient *announcement_api.ClientWithResponses,
	funchClient *funch_api.ClientWithResponses,
	userClient *user_api.ClientWithResponses,
) *Handler {
	if academicClient == nil {
		panic("academicClient is required")
	}
	if announcementClient == nil {
		panic("announcementClient is required")
	}
	if funchClient == nil {
		panic("funchClient is required")
	}
	if userClient == nil {
		panic("userClient is required")
	}
	return &Handler{
		academicClient:     academicClient,
		announcementClient: announcementClient,
		funchClient:        funchClient,
		userClient:         userClient,
	}
}

var _ api.ServerInterface = (*Handler)(nil)
