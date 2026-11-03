// Package auth は App API の認証（Firebase App Check と Firebase Authentication）を行う。
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"firebase.google.com/go/v4/auth"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
)

const (
	appCheckHeader      = "X-Firebase-AppCheck"
	authorizationHeader = "Authorization"
)

// AppCheckVerifier は App Check トークンを検証する。
type AppCheckVerifier interface {
	VerifyToken(token string) error
}

// IDTokenVerifier は Firebase ID トークンを検証し、ユーザーを返す。
type IDTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (User, error)
}

// User は認証済みのユーザー。
type User struct {
	ID    string
	Email string
}

// userRequiredProcedures はログイン中のユーザーを必要とする RPC。
var userRequiredProcedures = map[string]bool{
	appv1connect.CancelledClassServiceListCancelledClassesProcedure:            true,
	appv1connect.CourseRegistrationServiceListCourseRegistrationsProcedure:     true,
	appv1connect.CourseRegistrationServiceCreateCourseRegistrationProcedure:    true,
	appv1connect.CourseRegistrationServiceDeleteCourseRegistrationProcedure:    true,
	appv1connect.FcmTokenServiceRegisterFcmTokenProcedure:                      true,
	appv1connect.MakeupClassServiceListMakeupClassesProcedure:                  true,
	appv1connect.PersonalCalendarItemServiceListPersonalCalendarItemsProcedure: true,
	appv1connect.ReservationServiceListReservationsProcedure:                   true,
	appv1connect.RoomChangeServiceListRoomChangesProcedure:                     true,
	appv1connect.UserServiceGetUserProcedure:                                   true,
	appv1connect.UserServiceUpdateUserProcedure:                                true,
}

// RequiresUser は RPC がログイン中のユーザーを必要とするかを返す。
func RequiresUser(procedure string) bool {
	return userRequiredProcedures[procedure]
}

// NewInterceptor はすべての RPC で App Check トークンを検証し、
// ユーザーを必要とする RPC では ID トークンを検証してユーザーを context に格納する Interceptor を返す。
func NewInterceptor(appCheck AppCheckVerifier, idToken IDTokenVerifier) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if err := verifyAppCheck(appCheck, req.Header()); err != nil {
				return nil, err
			}
			if RequiresUser(req.Spec().Procedure) {
				user, err := verifyUser(ctx, idToken, req.Header())
				if err != nil {
					return nil, err
				}
				ctx = WithUser(ctx, user)
			}
			return next(ctx, req)
		}
	}
}

func verifyAppCheck(v AppCheckVerifier, h http.Header) error {
	token := strings.TrimPrefix(h.Get(appCheckHeader), "Bearer ")
	if token == "" {
		return unauthenticated("APP_CHECK_TOKEN_MISSING", appCheckHeader+" header is required")
	}
	if err := v.VerifyToken(token); err != nil {
		return unauthenticated("APP_CHECK_TOKEN_INVALID", appCheckHeader+" header value is invalid")
	}
	return nil
}

func verifyUser(ctx context.Context, v IDTokenVerifier, h http.Header) (User, error) {
	authorization := h.Get(authorizationHeader)
	token, ok := strings.CutPrefix(authorization, "Bearer ")
	if !ok || token == "" {
		return User{}, unauthenticated("ID_TOKEN_MISSING", "Authorization header with a Bearer ID token is required")
	}
	user, err := v.VerifyIDToken(ctx, token)
	if err != nil {
		return User{}, unauthenticated("ID_TOKEN_INVALID", "ID token is invalid or expired")
	}
	return user, nil
}

type userKey struct{}

// WithUser はユーザーを context に格納する。
func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// UserFromContext は context からユーザーを取り出す。
func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey{}).(User)
	return user, ok && user.ID != ""
}

func unauthenticated(reason, message string) error {
	return newError(connect.CodeUnauthenticated, reason, message)
}

// FirebaseIDTokenVerifier は Firebase Authentication のクライアントで ID トークンを検証する。
type FirebaseIDTokenVerifier struct {
	Client *auth.Client
}

// VerifyIDToken は ID トークンを検証する。
func (v FirebaseIDTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (User, error) {
	if v.Client == nil {
		return User{}, errors.New("auth client is not configured")
	}
	token, err := v.Client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return User{}, err
	}
	user := User{ID: token.UID}
	if email, ok := token.Claims["email"].(string); ok {
		user.Email = email
	}
	return user, nil
}

// AppCheckTokenVerifierFunc は関数を AppCheckVerifier として扱う。
type AppCheckTokenVerifierFunc func(token string) error

// VerifyToken は App Check トークンを検証する。
func (f AppCheckTokenVerifierFunc) VerifyToken(token string) error { return f(token) }
