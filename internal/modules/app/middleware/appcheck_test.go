package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppCheckMiddleware(t *testing.T) {
	t.Setenv("GIN_MODE", "test")
	gin.SetMode(gin.TestMode)

	t.Run("AppCheck が不要な場合はヘッダーがなくても通す", func(t *testing.T) {
		router := gin.New()
		router.Use(AppCheckMiddleware(nil))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("AppCheck が必要な場合はヘッダーが必須", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set(api.FirebaseAppCheckAuthScopes, []string{})
			c.Next()
		})
		router.Use(AppCheckMiddleware(nil))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"error":"X-Firebase-AppCheck header is required"}`, rec.Body.String())
	})
}
