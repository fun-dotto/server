package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// abortJSON は認証・認可エラーレスポンスを既存 BFF と揃った形式で返す。
// レスポンスボディは {"error": "<message>"} 固定。
func abortJSON(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}

// abortUnauthorized は 401 を返す。
func abortUnauthorized(c *gin.Context, message string) {
	abortJSON(c, http.StatusUnauthorized, message)
}

// abortForbidden は 403 を返す。
func abortForbidden(c *gin.Context, message string) {
	abortJSON(c, http.StatusForbidden, message)
}
