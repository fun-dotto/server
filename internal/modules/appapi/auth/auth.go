// Package auth は App API の認証（Firebase App Check と Firebase Authentication）を行う。
package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"connectrpc.com/connect"
	"firebase.google.com/go/v4/appcheck"
	"firebase.google.com/go/v4/auth"
	"github.com/fun-dotto/server/gen/app/v1/appv1connect"
	"github.com/fun-dotto/server/internal/shared/cfaccess"
)

const (
	appCheckHeader      = "X-Firebase-AppCheck"
	authorizationHeader = "Authorization"
	// serviceAuthHeader は Web などのサービスが Google の ID トークンを付与するヘッダー。
	// Authorization はユーザーの Firebase ID トークンに使うため別のヘッダーにする。
	serviceAuthHeader = "X-Service-Authorization"
)

// AppCheckVerifier は App Check トークンを検証し、使用済みにする（*appcheck.Client が満たす）。
type AppCheckVerifier interface {
	VerifyOneTimeToken(ctx context.Context, token string) (*appcheck.DecodedAppCheckToken, error)
}

// AppCheck は App Check トークンを使い捨てとして検証し、許可されたアプリからのものかを確かめる。
type AppCheck struct {
	verifier      AppCheckVerifier
	allowedAppIDs []string
}

// NewAppCheck は Verifier と、カンマ区切りの許可する Firebase アプリ ID から AppCheck を作る。
func NewAppCheck(verifier AppCheckVerifier, allowedAppIDs string) (*AppCheck, error) {
	if verifier == nil {
		return nil, errors.New("app check verifier is nil")
	}
	var ids []string
	for id := range strings.SplitSeq(allowedAppIDs, ",") {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil, errors.New("app check allowed app ids is empty")
	}
	return &AppCheck{verifier: verifier, allowedAppIDs: ids}, nil
}

// CFAccessVerifier は Cloudflare Access の JWT を検証する。
type CFAccessVerifier interface {
	Verify(ctx context.Context, token string) error
}

// ServiceVerifier は呼び出し元サービスの Google ID トークンを検証する。
type ServiceVerifier interface {
	Verify(ctx context.Context, token string) error
}

// ErrIDTokenRevoked は ID トークンが失効している（ログアウトやパスワード変更など）ことを表す。
var ErrIDTokenRevoked = errors.New("id token has been revoked")

// ErrUserDisabled はユーザーが無効化されていることを表す。
var ErrUserDisabled = errors.New("user has been disabled")

// IDTokenVerifier は Firebase ID トークンを検証し、ユーザーを返す。
// 失効したトークンには ErrIDTokenRevoked、無効化されたユーザーには ErrUserDisabled を返す。
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

// NewInterceptor はすべての RPC で App Check トークンを使い捨てとして検証し、
// ユーザーを必要とする RPC では ID トークンを検証してユーザーを context に格納する Interceptor を返す。
// 有効な Cloudflare Access の JWT、または許可されたサービスの Google ID トークンを持つリクエストは
// App Check の検証を省略する。
func NewInterceptor(appCheck *AppCheck, cfAccess CFAccessVerifier, service ServiceVerifier, idToken IDTokenVerifier) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if !verifyService(ctx, service, req.Header()) && !verifyCFAccess(ctx, cfAccess, req.Header()) {
				if err := appCheck.verify(ctx, req.Header()); err != nil {
					return nil, err
				}
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

// verifyCFAccess は Cloudflare Access の JWT が有効かを返す。JWT がない、または不正な場合は false を返す。
// v が nil（Cloudflare Access を使わない環境）の場合は常に false を返す。
func verifyCFAccess(ctx context.Context, v CFAccessVerifier, h http.Header) bool {
	token := h.Get(cfaccess.JWTHeader)
	if v == nil || token == "" {
		return false
	}
	if err := v.Verify(ctx, token); err != nil {
		slog.WarnContext(ctx, "cloudflare access jwt verification failed", "error", err)
		return false
	}
	return true
}

// verifyService は呼び出し元サービスの Google ID トークンが有効かを返す。トークンがない、または不正な場合は false を返す。
func verifyService(ctx context.Context, v ServiceVerifier, h http.Header) bool {
	token, ok := strings.CutPrefix(h.Get(serviceAuthHeader), "Bearer ")
	if !ok || token == "" {
		return false
	}
	if err := v.Verify(ctx, token); err != nil {
		slog.WarnContext(ctx, "service id token verification failed", "error", err)
		return false
	}
	return true
}

// verify は App Check トークンを検証して使用済みにし、許可されたアプリからのものかを確かめる。
// 使用済みのトークン（リプレイ）は拒否する。
func (a *AppCheck) verify(ctx context.Context, h http.Header) error {
	token := strings.TrimPrefix(h.Get(appCheckHeader), "Bearer ")
	if token == "" {
		return unauthenticated("APP_CHECK_TOKEN_MISSING", appCheckHeader+" header is required")
	}
	decoded, err := a.verifier.VerifyOneTimeToken(ctx, token)
	if err != nil {
		slog.WarnContext(ctx, "app check token verification failed", "error", err)
		return unauthenticated("APP_CHECK_TOKEN_INVALID", appCheckHeader+" header value is invalid")
	}
	if decoded.AlreadyConsumed == nil || *decoded.AlreadyConsumed {
		return unauthenticated("APP_CHECK_TOKEN_ALREADY_CONSUMED", appCheckHeader+" header value has already been used")
	}
	if !slices.Contains(a.allowedAppIDs, decoded.AppID) {
		slog.WarnContext(ctx, "app check token from disallowed app", "app_id", decoded.AppID)
		return unauthenticated("APP_CHECK_APP_NOT_ALLOWED", appCheckHeader+" header value is not issued for an allowed app")
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
	switch {
	case errors.Is(err, ErrIDTokenRevoked):
		return User{}, unauthenticated("ID_TOKEN_REVOKED", "ID token has been revoked")
	case errors.Is(err, ErrUserDisabled):
		return User{}, unauthenticated("USER_DISABLED", "user has been disabled")
	case err != nil:
		slog.WarnContext(ctx, "id token verification failed", "error", err)
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

// VerifyIDToken は ID トークンを検証し、トークンの失効とユーザーの無効化も確かめる。
// 失効・無効化の確認のため、毎回 Firebase Authentication に問い合わせる。
func (v FirebaseIDTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (User, error) {
	if v.Client == nil {
		return User{}, errors.New("auth client is not configured")
	}
	token, err := v.Client.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	switch {
	case auth.IsIDTokenRevoked(err):
		return User{}, ErrIDTokenRevoked
	case auth.IsUserDisabled(err):
		return User{}, ErrUserDisabled
	case err != nil:
		return User{}, err
	}
	user := User{ID: token.UID}
	if email, ok := token.Claims["email"].(string); ok {
		user.Email = email
	}
	return user, nil
}
