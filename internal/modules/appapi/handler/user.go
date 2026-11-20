package handler

import (
	"context"

	"connectrpc.com/connect"
	appv1 "github.com/fun-dotto/server/gen/app/v1"
	"github.com/fun-dotto/server/internal/modules/appapi/auth"
	"github.com/fun-dotto/server/internal/shared/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetUser はログイン中のユーザーを返す。
func (h *Handler) GetUser(ctx context.Context, _ *connect.Request[appv1.GetUserRequest]) (*connect.Response[appv1.GetUserResponse], error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var record model.User
	if err := h.db.WithContext(ctx).First(&record, "id = ?", userID).Error; err != nil {
		if isNotFound(err) {
			return nil, newError(connect.CodeNotFound, "USER_NOT_FOUND", "user not found; create the user with UpdateUser", nil)
		}
		return nil, internalError(ctx, "GetUser", err)
	}
	return connect.NewResponse(&appv1.GetUserResponse{User: toProtoUser(record)}), nil
}

// UpdateUser はログイン中のユーザーを更新する。存在しない場合は作成する。
func (h *Handler) UpdateUser(ctx context.Context, req *connect.Request[appv1.UpdateUserRequest]) (*connect.Response[appv1.UpdateUserResponse], error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, unauthenticated()
	}
	m := req.Msg

	// update_mask が空の場合は値が入っているフィールドをすべて更新する。
	paths := map[string]bool{"grade": m.Grade != nil, "course": m.Course != nil, "class": m.Class != nil}
	var v violations
	if mask := m.GetUpdateMask().GetPaths(); len(mask) > 0 {
		for k := range paths {
			paths[k] = false
		}
		for _, p := range mask {
			if _, ok := paths[p]; !ok {
				v.add("update_mask", "contains an unknown field: "+p)
				continue
			}
			paths[p] = true
		}
	}
	grade := optionalFromProto(gradeMap, m.Grade, "grade", &v)
	course := optionalFromProto(courseMap, m.Course, "course", &v)
	class := optionalFromProto(classMap, m.Class, "class", &v)
	if err := v.err(); err != nil {
		return nil, err
	}

	var record model.User
	err := h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&record, "id = ?", user.ID).Error
		if err != nil && !isNotFound(err) {
			return err
		}
		record.ID = user.ID
		if user.Email != "" {
			record.Email = user.Email
		}
		if paths["grade"] {
			record.Grade = grade
		}
		if paths["course"] {
			record.Course = course
		}
		if paths["class"] {
			record.Class = class
		}
		return tx.Save(&record).Error
	})
	if err != nil {
		return nil, internalError(ctx, "UpdateUser", err)
	}
	return connect.NewResponse(&appv1.UpdateUserResponse{User: toProtoUser(record)}), nil
}

// optionalFromProto は optional な enum を DB の文字列に変換する。未設定・UNSPECIFIED の場合は nil（値を消す）。
func optionalFromProto[P comparable](m enumMap[P], p *P, field string, v *violations) *string {
	if p == nil {
		return nil
	}
	var zero P
	if *p == zero {
		return nil
	}
	s, ok := m.fromProto[*p]
	if !ok {
		v.add(field, "has an unknown value")
		return nil
	}
	return &s
}

func toProtoUser(u model.User) *appv1.User {
	return &appv1.User{
		Grade:  gradeMap.toProtoPtr(u.Grade),
		Course: courseMap.toProtoPtr(u.Course),
		Class:  classMap.toProtoPtr(u.Class),
	}
}

// RegisterFcmToken はログイン中のユーザーの FCM トークンを登録する。
// 存在しない場合は作成し、存在する場合は更新日時を更新する。他のユーザーのトークンの場合は PERMISSION_DENIED を返す。
func (h *Handler) RegisterFcmToken(ctx context.Context, req *connect.Request[appv1.RegisterFcmTokenRequest]) (*connect.Response[appv1.RegisterFcmTokenResponse], error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	token := req.Msg.GetToken()
	if token == "" {
		return nil, invalidArgument("token", "is required")
	}

	var record model.FCMToken
	err = h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&record, "token = ?", token).Error
		switch {
		case err == nil:
			if record.UserID != userID {
				return permissionDenied("fcm_token", nil)
			}
			return tx.Model(&record).Update("updated_at", h.now()).Error
		case isNotFound(err):
			record = model.FCMToken{Token: token, UserID: userID}
			return tx.Omit("User").Create(&record).Error
		default:
			return err
		}
	})
	if err != nil {
		if connect.CodeOf(err) == connect.CodePermissionDenied {
			return nil, err
		}
		if postgresErrorCode(err) == pgForeignKeyViolation {
			return nil, userNotFound()
		}
		return nil, internalError(ctx, "RegisterFcmToken", err)
	}
	return connect.NewResponse(&appv1.RegisterFcmTokenResponse{FcmToken: &appv1.FcmToken{
		Token:     record.Token,
		CreatedAt: toTimestamp(record.CreatedAt),
		UpdatedAt: toTimestamp(record.UpdatedAt),
	}}), nil
}
