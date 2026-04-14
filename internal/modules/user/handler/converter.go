package handler

import (
	api "github.com/fun-dotto/user-api/generated"
	"github.com/fun-dotto/user-api/internal/domain"
)

func toAPIUser(u domain.User) api.User {
	user := api.User{
		Id:    u.ID,
		Email: u.Email,
	}

	if u.Grade != nil {
		grade := api.DottoFoundationV1Grade(*u.Grade)
		user.Grade = &grade
	}
	if u.Course != nil {
		course := api.DottoFoundationV1Course(*u.Course)
		user.Course = &course
	}
	if u.Class != nil {
		class := api.DottoFoundationV1Class(*u.Class)
		user.Class = &class
	}

	return user
}

func toDomainUser(id string, req api.UserRequest) domain.User {
	user := domain.User{
		ID:    id,
		Email: req.Email,
	}

	if req.Grade != nil {
		grade := string(*req.Grade)
		user.Grade = &grade
	}
	if req.Course != nil {
		course := string(*req.Course)
		user.Course = &course
	}
	if req.Class != nil {
		class := string(*req.Class)
		user.Class = &class
	}

	return user
}

func toAPIUsers(users []domain.User) []api.User {
	apiUsers := make([]api.User, 0, len(users))
	for _, u := range users {
		apiUsers = append(apiUsers, toAPIUser(u))
	}
	return apiUsers
}

func toAPIFCMToken(t domain.FCMToken) api.FCMToken {
	return api.FCMToken{
		Token:     t.Token,
		UserId:    t.UserID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toAPIFCMTokens(tokens []domain.FCMToken) []api.FCMToken {
	apiTokens := make([]api.FCMToken, 0, len(tokens))
	for _, t := range tokens {
		apiTokens = append(apiTokens, toAPIFCMToken(t))
	}
	return apiTokens
}

func toDomainFCMToken(req api.FCMTokenRequest) domain.FCMToken {
	return domain.FCMToken{
		Token:  req.Token,
		UserID: req.UserId,
	}
}

func toAPINotification(n domain.Notification) api.Notification {
	return api.Notification{
		Id:            n.ID,
		Title:         n.Title,
		Message:       n.Message,
		Url:           n.URL,
		NotifyAt:      n.NotifyAt,
		IsNotified:    n.IsNotified,
		TargetUserIds: n.TargetUserIDs,
	}
}

func toAPINotifications(notifications []domain.Notification) []api.Notification {
	result := make([]api.Notification, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toAPINotification(n))
	}
	return result
}

func toDomainNotification(id string, req api.NotificationRequest) domain.Notification {
	return domain.Notification{
		ID:            id,
		Title:         req.Title,
		Message:       req.Message,
		URL:           req.Url,
		NotifyAt:      req.NotifyAt,
		TargetUserIDs: req.TargetUserIds,
	}
}

func toDomainNotificationListFilter(params api.NotificationV1ListParams) domain.NotificationListFilter {
	return domain.NotificationListFilter{
		NotifyAtFrom: params.NotifyAtFrom,
		NotifyAtTo:   params.NotifyAtTo,
		IsNotified:   params.IsNotified,
	}
}

func toDomainFCMTokenListFilter(params api.FCMTokenV1ListParams) domain.FCMTokenListFilter {
	filter := domain.FCMTokenListFilter{
		UpdatedAtFrom: params.UpdatedAtFrom,
		UpdatedAtTo:   params.UpdatedAtTo,
	}

	if params.UserIds != nil {
		filter.UserIDs = append(filter.UserIDs, (*params.UserIds)...)
	}
	if params.Tokens != nil {
		filter.Tokens = append(filter.Tokens, (*params.Tokens)...)
	}

	return filter
}
