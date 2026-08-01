package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/server/gen/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	t.Setenv("GIN_MODE", "test")
	gin.SetMode(gin.TestMode)

	t.Run("BearerAuth が不要な場合は Authorization がなくても通す", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthMiddleware(nil))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("BearerAuth が必要な場合は Authorization が必須", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set(api.BearerAuthScopes, []string{})
			c.Next()
		})
		router.Use(AuthMiddleware(nil))
		router.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.JSONEq(t, `{"error":"Authorization header is required"}`, rec.Body.String())
	})
}

func TestUserIDFromContext(t *testing.T) {
	t.Run("context に userID が存在する", func(t *testing.T) {
		ctx := contextWithUserID("user-123")
		got, ok := UserIDFromContext(ctx)
		require.True(t, ok)
		require.Equal(t, "user-123", got)
	})

	t.Run("context に userID が存在しない", func(t *testing.T) {
		got, ok := UserIDFromContext(t.Context())
		require.False(t, ok)
		require.Empty(t, got)
	})
}

func TestUserEmailFromContext(t *testing.T) {
	t.Run("context に userEmail が存在する", func(t *testing.T) {
		ctx := contextWithUserEmail("user@example.com")
		got, ok := UserEmailFromContext(ctx)
		require.True(t, ok)
		require.Equal(t, "user@example.com", got)
	})

	t.Run("context に userEmail が存在しない", func(t *testing.T) {
		got, ok := UserEmailFromContext(t.Context())
		require.False(t, ok)
		require.Empty(t, got)
	})
}

func contextWithUserID(userID string) context.Context {
	return context.WithValue(context.Background(), userIDKey, userID)
}

func contextWithUserEmail(email string) context.Context {
	return context.WithValue(context.Background(), userEmailKey, email)
}
