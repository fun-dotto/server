package auth

import "github.com/gin-gonic/gin"

// 既知のカスタムクレーム名。admin-bff-api 側との互換性のためそのまま採用する。
const (
	ClaimAdmin     = "admin"
	ClaimDeveloper = "developer"
)

// RequireAdmin は admin / developer のいずれかのカスタムクレームを持つ Bearer
// での呼び出しのみを許可する。満たさない場合は 403 を返し、ハンドラは
// そのまま return すれば良い (ok=false 時)。
//
// 前段で RequireUserUnlessAllowed を通過しているため Bearer 検証は済んでいる前提。
// 想定外に未認証で到達した場合は 401 を返す。
func RequireAdmin(c *gin.Context) bool {
	ctx := c.Request.Context()

	if _, ok := GetFirebaseToken(ctx); !ok {
		abortUnauthorized(c, "authentication required")
		return false
	}

	if HasClaim(ctx, ClaimAdmin) || HasClaim(ctx, ClaimDeveloper) {
		return true
	}

	abortForbidden(c, "insufficient permissions")
	return false
}

// RequireAnyClaim は指定したカスタムクレームのいずれかが boolean true で
// 設定されていれば true を返す。それ以外は 403 を返し false を返す。
func RequireAnyClaim(c *gin.Context, claims ...string) bool {
	ctx := c.Request.Context()

	if _, ok := GetFirebaseToken(ctx); !ok {
		abortUnauthorized(c, "authentication required")
		return false
	}

	for _, claim := range claims {
		if HasClaim(ctx, claim) {
			return true
		}
	}

	abortForbidden(c, "insufficient permissions")
	return false
}
