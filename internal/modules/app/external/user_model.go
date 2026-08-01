package external

import (
	user_api "github.com/fun-dotto/server/gen/user"
	"github.com/fun-dotto/server/internal/modules/app/domain"
)

// ToDomainUser は外部APIのUserをDomainのUserに変換する
func ToDomainUser(m user_api.User) domain.User {
	user := domain.User{
		ID:    m.Id,
		Email: m.Email,
	}
	if m.Grade != nil {
		g := domain.Grade(*m.Grade)
		user.Grade = &g
	}
	if m.Course != nil {
		c := domain.Course(*m.Course)
		user.Course = &c
	}
	if m.Class != nil {
		cl := domain.Class(*m.Class)
		user.Class = &cl
	}
	return user
}

// ToExternalUserRequest はDomainのUserRequestを外部APIのUserRequestに変換する
func ToExternalUserRequest(req domain.UserRequest) user_api.UserRequest {
	body := user_api.UserRequest{
		Email: req.Email,
	}
	if req.Grade != nil {
		g := user_api.DottoFoundationV1Grade(*req.Grade)
		body.Grade = &g
	}
	if req.Course != nil {
		c := user_api.DottoFoundationV1Course(*req.Course)
		body.Course = &c
	}
	if req.Class != nil {
		cl := user_api.DottoFoundationV1Class(*req.Class)
		body.Class = &cl
	}
	return body
}

func ToDomainFCMToken(m user_api.FCMToken) domain.FCMToken {
	return domain.FCMToken{
		Token:     m.Token,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToExternalFCMTokenRequest(userID string, req domain.FCMTokenRequest) user_api.FCMTokenRequest {
	return user_api.FCMTokenRequest{
		Token:  req.Token,
		UserId: userID,
	}
}
