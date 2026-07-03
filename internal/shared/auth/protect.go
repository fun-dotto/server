package auth

import (
	"github.com/gin-gonic/gin"
)

// AllowList は「ログインしていなくても叩いてよい」エンドポイントの集合。
// キーは "<METHOD> <Gin の登録パス>" 形式。例: "GET /v1/announcements".
// 登録パスは gin.Context.FullPath() の値と一致させる必要がある
// （path パラメータは ":id" のような Gin 形式で書く）。
type AllowList map[string]struct{}

// NewAllowList は文字列スライスから AllowList を構築する小ヘルパ。
func NewAllowList(routes ...string) AllowList {
	a := make(AllowList, len(routes))
	for _, r := range routes {
		a[r] = struct{}{}
	}
	return a
}

// RequireUserUnlessAllowed は default-deny で「ログイン必須」を強制する middleware。
// allowList に登録された経路では Bearer 認証が無くても通すが、
// その場合でも AppCheck が検証成功している必要がある（=正規アプリ確認）。
// allowList に登録されていない経路は Bearer 検証成功が必須。
//
// 前段に Extract を必ず登録しておくこと。
func RequireUserUnlessAllowed(allowList AllowList) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		key := c.Request.Method + " " + c.FullPath()

		if _, allowed := allowList[key]; allowed {
			// 匿名アクセスを許可する経路。Bearer or AppCheck のいずれかが
			// 検証成功していれば通す。
			if _, ok := GetFirebaseToken(ctx); ok {
				c.Next()
				return
			}
			if IsAppCheckVerified(ctx) {
				c.Next()
				return
			}
			abortUnauthorized(c, "authentication required")
			return
		}

		// 通常経路。Bearer 検証成功必須。
		if _, ok := GetFirebaseToken(ctx); !ok {
			abortUnauthorized(c, "authentication required")
			return
		}
		c.Next()
	}
}
