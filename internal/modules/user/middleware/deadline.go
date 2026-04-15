package middleware

import (
	"context"
	"errors"
	"net/http"

	api "github.com/fun-dotto/user-api/generated"
	"github.com/gin-gonic/gin"
)

// DeadlineErrorMapper は context.DeadlineExceeded を 504 にマップする strict middleware。
// gin-contrib/timeout が先に 504 を書き込んでいる場合、この書き込みは gin に破棄される。
// いずれの経路でも 504 を返す動作は同じであり、二重防御として機能する。
func DeadlineErrorMapper() api.StrictMiddlewareFunc {
	return func(f api.StrictHandlerFunc, _ string) api.StrictHandlerFunc {
		return func(c *gin.Context, request any) (any, error) {
			resp, err := f(c, request)
			if err != nil && errors.Is(err, context.DeadlineExceeded) {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
					"msg": "request timed out",
				})
				return nil, nil
			}
			return resp, err
		}
	}
}
