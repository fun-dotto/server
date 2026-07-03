package auth

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

// Gin / context.Context に検証済みの認証情報を出し入れするためのキー。
// 衝突を避けるため private な struct を採用する。
type (
	firebaseTokenKey struct{}
	appCheckOKKey    struct{}
)

var (
	tokenCtxKey    = firebaseTokenKey{}
	appCheckCtxKey = appCheckOKKey{}
)

// setFirebaseToken は検証済みの Firebase ID Token を Gin / context の両方に格納する。
// 後段のハンドラからは GetFirebaseToken / UserID を通じて参照する。
func setFirebaseToken(c *gin.Context, token *auth.Token) {
	ctx := context.WithValue(c.Request.Context(), tokenCtxKey, token)
	c.Request = c.Request.WithContext(ctx)
	c.Set(keyOf(tokenCtxKey), token)
}

// markAppCheckOK は AppCheck トークンの検証が成功したことを context に記録する。
func markAppCheckOK(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), appCheckCtxKey, true)
	c.Request = c.Request.WithContext(ctx)
	c.Set(keyOf(appCheckCtxKey), true)
}

// GetFirebaseToken は context.Context から検証済み Firebase Token を取り出す。
// Extract ミドルウェアを通過していない、もしくは検証に失敗している場合は ok=false。
func GetFirebaseToken(ctx context.Context) (*auth.Token, bool) {
	v := ctx.Value(tokenCtxKey)
	if v == nil {
		return nil, false
	}
	token, ok := v.(*auth.Token)
	return token, ok
}

// UserID は検証済みトークンの UID を返す。未認証の場合は空文字列。
func UserID(ctx context.Context) string {
	token, ok := GetFirebaseToken(ctx)
	if !ok || token == nil {
		return ""
	}
	return token.UID
}

// IsAppCheckVerified は AppCheck の検証が成功済みかを返す。
func IsAppCheckVerified(ctx context.Context) bool {
	v := ctx.Value(appCheckCtxKey)
	if v == nil {
		return false
	}
	ok, _ := v.(bool)
	return ok
}

// HasClaim は検証済み Firebase Token に対してカスタムクレーム名 claim が
// boolean true で設定されているかを返す。
func HasClaim(ctx context.Context, claim string) bool {
	token, ok := GetFirebaseToken(ctx)
	if !ok || token == nil {
		return false
	}
	v, exists := token.Claims[claim]
	if !exists {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

// keyOf は Gin の c.Set 用に struct キーを文字列化する。
// Gin の context は string キーしか扱わないため、リフレクションを避けて固定文字列に寄せる。
func keyOf(key any) string {
	switch key.(type) {
	case firebaseTokenKey:
		return "auth.firebaseToken"
	case appCheckOKKey:
		return "auth.appCheckOK"
	default:
		return ""
	}
}
