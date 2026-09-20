package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"
)

type firebaseTokenKey struct{}
type authErrorStatusKey struct{}
type authErrorMessageKey struct{}

// FirebaseTokenContextKey は Gin の context および context.Context に
// Firebase ID トークンの検証結果を格納するキーです。
var FirebaseTokenContextKey = firebaseTokenKey{}

var (
	authenticationErrorStatusKey  = authErrorStatusKey{}
	authenticationErrorMessageKey = authErrorMessageKey{}
)

type AuthenticationError struct {
	StatusCode int
	Message    string
}

func (e *AuthenticationError) Error() string {
	return e.Message
}

// FirebaseAuthenticationFunc は OpenAPI validator 向けの AuthenticationFunc を返します。
// 認証に成功した場合は検証済みトークンを Gin / request context に格納します。
func FirebaseAuthenticationFunc(authClient *auth.Client) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, _ *openapi3filter.AuthenticationInput) error {
		ginCtx := ginmiddleware.GetGinContext(ctx)
		if ginCtx == nil {
			return &AuthenticationError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Authentication context is unavailable",
			}
		}

		token, err := verifyFirebaseToken(ginCtx.GetHeader("Authorization"), ginCtx.Request.Context(), authClient)
		if err != nil {
			var authErr *AuthenticationError
			if errors.As(err, &authErr) {
				ginCtx.Set(authenticationErrorStatusKey, authErr.StatusCode)
				ginCtx.Set(authenticationErrorMessageKey, authErr.Message)
			}
			return err
		}

		setFirebaseToken(ginCtx, token)
		return nil
	}
}

// FirebaseAuth は Authorization: Bearer <Firebase ID Token> を検証する Gin ミドルウェアです。
// 検証に成功すると、デコードされたトークン（*auth.Token）を context に格納して次のハンドラに渡します。
func FirebaseAuth(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := verifyFirebaseToken(c.GetHeader("Authorization"), c.Request.Context(), authClient)
		if err != nil {
			var authErr *AuthenticationError
			if errors.As(err, &authErr) {
				c.AbortWithStatusJSON(authErr.StatusCode, gin.H{"error": authErr.Message})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}

		setFirebaseToken(c, token)
		c.Next()
	}
}

func verifyFirebaseToken(authHeader string, ctx context.Context, authClient *auth.Client) (*auth.Token, error) {
	if authHeader == "" {
		return nil, &AuthenticationError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Authorization header is required",
		}
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, &AuthenticationError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Authorization header must be Bearer <token>",
		}
	}

	idToken := strings.TrimSpace(parts[1])
	if idToken == "" {
		return nil, &AuthenticationError{
			StatusCode: http.StatusUnauthorized,
			Message:    "ID token is required",
		}
	}

	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		status, message := authErrorResponse(err)
		return nil, &AuthenticationError{
			StatusCode: status,
			Message:    message,
		}
	}

	return token, nil
}

func setFirebaseToken(c *gin.Context, token *auth.Token) {
	ctx := context.WithValue(c.Request.Context(), FirebaseTokenContextKey, token)
	c.Set(FirebaseTokenContextKey, token)
	c.Request = c.Request.WithContext(ctx)
}

// GetAuthenticationError は validator の AuthenticationFunc が格納した認証失敗情報を返します。
func GetAuthenticationError(c *gin.Context) (int, string, bool) {
	status, statusExists := c.Get(authenticationErrorStatusKey)
	message, messageExists := c.Get(authenticationErrorMessageKey)
	if !statusExists || !messageExists {
		return 0, "", false
	}

	statusCode, ok := status.(int)
	if !ok {
		return 0, "", false
	}
	errorMessage, ok := message.(string)
	if !ok {
		return 0, "", false
	}

	return statusCode, errorMessage, true
}

// authErrorResponse は Firebase Auth の検証エラーから HTTP ステータスとメッセージを返します。
func authErrorResponse(err error) (int, string) {
	if auth.IsIDTokenExpired(err) {
		return http.StatusUnauthorized, "ID token has expired"
	}
	if auth.IsIDTokenInvalid(err) {
		return http.StatusUnauthorized, "Invalid ID token"
	}
	if auth.IsIDTokenRevoked(err) {
		return http.StatusUnauthorized, "ID token has been revoked"
	}
	if auth.IsUserDisabled(err) {
		return http.StatusForbidden, "User has been disabled"
	}
	return http.StatusUnauthorized, "Authentication failed"
}

// GetFirebaseToken は Gin の context から検証済みの Firebase トークンを取得します。
// FirebaseAuth ミドルウェア通過後のハンドラ内でのみ有効です。
func GetFirebaseToken(c *gin.Context) (*auth.Token, bool) {
	val, exists := c.Get(FirebaseTokenContextKey)
	if !exists {
		return nil, false
	}
	token, ok := val.(*auth.Token)
	return token, ok
}

// GetFirebaseUID は Gin の context から認証済みユーザーの UID を取得します。
// トークンが存在しない場合は空文字列を返します。
func GetFirebaseUID(c *gin.Context) string {
	token, ok := GetFirebaseToken(c)
	if !ok || token == nil {
		return ""
	}
	return token.UID
}

// GetFirebaseTokenFromContext は context.Context から Firebase トークンを取得します。
// ミドルウェアで c.Request に設定した context から取得する場合に使用します。
func GetFirebaseTokenFromContext(ctx context.Context) (*auth.Token, bool) {
	val := ctx.Value(FirebaseTokenContextKey)
	if val == nil {
		return nil, false
	}
	token, ok := val.(*auth.Token)
	return token, ok
}

// RequireAnyClaim は指定したカスタムクレームのいずれかが true であることを検証します。
// いずれのクレームも満たさない場合は 403 レスポンスを返し false を返します。
func RequireAnyClaim(c *gin.Context, claims ...string) bool {
	token, ok := GetFirebaseToken(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication required",
		})
		return false
	}

	for _, claim := range claims {
		if val, exists := token.Claims[claim]; exists {
			if boolVal, ok := val.(bool); ok && boolVal {
				return true
			}
		}
	}

	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error": "Insufficient permissions",
	})
	return false
}
