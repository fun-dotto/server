package auth

import (
	"context"
	"log"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	appCheckHeader      = "X-Firebase-AppCheck"
)

// Extract は全リクエストに対して best-effort で認証情報を検証する pre middleware。
// 検証に失敗・トークン欠落の場合でも 401 にはせず、後段の RequireUserUnlessAllowed
// と各ハンドラの gate ヘルパに委ねる。
//
// 動作概要:
//   - Authorization: Bearer ヘッダがあれば失効チェック付きで Firebase ID Token を検証。
//     成功時は context に *auth.Token を格納する。
//   - Bearer 検証が成功した場合、AppCheck の検証はスキップする
//     （ログイン済みユーザは AppCheck を送らなくても叩ける方針）。
//   - Bearer が無い／失敗した場合に限り、X-Firebase-AppCheck ヘッダがあれば検証。
//     成功時は context に検証済みフラグを記録する。
func Extract(clients *Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		bearerOK := tryVerifyBearer(c, ctx, clients.Auth)
		if !bearerOK {
			tryVerifyAppCheck(c, clients)
		}

		c.Next()
	}
}

// tryVerifyBearer は Authorization ヘッダから Bearer トークンを抽出して検証する。
// 検証失敗時はログを残し false を返す。ヘッダ欠落も false を返す。
func tryVerifyBearer(c *gin.Context, ctx context.Context, authClient *auth.Client) bool {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return false
	}
	if !strings.HasPrefix(header, bearerPrefix) {
		return false
	}
	idToken := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
	if idToken == "" {
		return false
	}

	token, err := authClient.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		// 検証失敗は best-effort なのでここでは握りつぶし、後段の gate に判断を委ねる。
		// デバッグ容易性のため log には残しておく（PII を避け UID は出さない）。
		log.Printf("auth: bearer token verification failed: %v", err)
		return false
	}

	setFirebaseToken(c, token)
	return true
}

// tryVerifyAppCheck は X-Firebase-AppCheck ヘッダがあれば検証する。
// AppCheck クライアント未注入の場合は何もしない。
func tryVerifyAppCheck(c *gin.Context, clients *Clients) {
	if clients.AppCheck == nil {
		return
	}
	header := c.GetHeader(appCheckHeader)
	if header == "" {
		return
	}
	// クライアントによっては "Bearer <token>" 形式で送ってくるため両対応する。
	token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
	if token == "" {
		return
	}

	if _, err := clients.AppCheck.VerifyToken(token); err != nil {
		log.Printf("auth: app check token verification failed: %v", err)
		return
	}
	markAppCheckOK(c)
}
