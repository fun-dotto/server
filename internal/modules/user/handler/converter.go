package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
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
